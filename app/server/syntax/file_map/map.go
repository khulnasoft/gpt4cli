package file_map

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"gpt4cli-server/syntax"
	"strings"

	"github.com/khulnasoft/gpt4cli/shared"
	tree_sitter "github.com/smacker/go-tree-sitter"
)

var verboseLogging = os.Getenv("VERBOSE_LOGGING") == "true"

// FileMap represents a file's important definitions
type FileMap struct {
	Definitions []Definition
}

type Definition struct {
	Type      string       // "function", "class", "key", "selector", "instruction" etc
	Signature string       // The full signature/header without implementation
	Comments  []string     // Any comments that precede this definition
	TagAttrs  []string     // For xml style markup tags, the class and id attributes
	TagReps   int          // For tags, the number of times this tag is repeated
	Line      int          // Line number where definition starts
	Children  []Definition // For parent types that can contain nested definitions
}

type Node struct {
	Type   string
	Lang   shared.TreeSitterLanguage
	TsNode *tree_sitter.Node
	Bytes  []byte
}

func MapFile(ctx context.Context, filename string, content []byte) (*FileMap, error) {
	// Get appropriate parser
	var parser *tree_sitter.Parser
	var lang shared.TreeSitterLanguage
	file := filepath.Base(filename)
	if strings.Contains(strings.ToLower(file), "dockerfile") {
		lang = shared.TreeSitterLanguageDockerfile
		parser = syntax.GetParserForLanguage(lang)

		if parser == nil {
			return nil, fmt.Errorf("no parser found for dockerfile")
		}
	} else {
		ext := filepath.Ext(filename)
		parser, lang, _, _ = syntax.GetParserForExt(ext)

		if parser == nil {
			return nil, fmt.Errorf("unsupported file type: %s", ext)
		}
	}

	// Parse file
	tree, err := parser.ParseCtx(ctx, nil, content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}
	defer tree.Close()

	// Create map
	m := &FileMap{
		Definitions: mapNode(tree.RootNode(), content, lang),
	}
	return m, nil
}

func mapNode(node *tree_sitter.Node, content []byte, lang shared.TreeSitterLanguage) []Definition {
	switch lang {
	case shared.TreeSitterLanguageHtml:
		return mapMarkup(content)
	case shared.TreeSitterLanguageSvelte:
		return mapSvelte(content)
	default:
		return mapTraditional(Node{
			Lang:   lang,
			TsNode: node,
			Bytes:  content,
		}, nil)
	}
}

// For traditional programming languages
func mapTraditional(baseNode Node, parentNode *Node) []Definition {
	var defs []Definition
	cursor := tree_sitter.NewTreeCursor(baseNode.TsNode)
	defer cursor.Close()

	if cursor.GoToFirstChild() {
		for {
			tsNode := cursor.CurrentNode()
			node := Node{
				Type:   tsNode.Type(),
				Lang:   baseNode.Lang,
				TsNode: tsNode,
				Bytes:  baseNode.Bytes,
			}

			if isIncludeAndContinueNode(node) {
				if verboseLogging {
					fmt.Println("include and continue node", cursor.CurrentNode().Type())
				}
				if !cursor.GoToNextSibling() {
					break
				}
				continue
			}

			if verboseLogging {
				fmt.Println()
				fmt.Println("node", node.Type)
				// fmt.Println("content", string(node.Content(content)))
				fmt.Println()
			}

			// Check if this is a definition node
			if isDefinitionNode(node, parentNode) {
				if verboseLogging {
					fmt.Println("definition node", node.Type)
				}

				def := Definition{
					Type: node.Type,
					Line: int(tsNode.StartPoint().Row) + 1,
				}

				if isAssignmentNode(node) {
					if verboseLogging {
						fmt.Println("assignment node", node.Type)
					}
					// Try different field names for identifiers
					// fmt.Printf("assignment node: %s\n", node.Type)
					sig := ""

					assignmentBoundary := findAssignmentBoundary(node)

					if assignmentBoundary != nil {
						start := tsNode.StartByte()
						end := assignmentBoundary.TsNode.StartByte()
						sig = string(node.Bytes[start:end])
						sig = strings.TrimSuffix(strings.TrimSpace(sig), "=")
					} else {
						identifiers := findIdentifier(node)
						if len(identifiers) > 0 {
							if verboseLogging {
								fmt.Println("found identifiers", len(identifiers))
							}

							start := tsNode.StartByte()
							end := identifiers[len(identifiers)-1].TsNode.EndByte()
							sig = string(node.Bytes[start:end])
						} else {
							if verboseLogging {
								fmt.Println("no identifier found", node.Type)
							}
							sig = string(node.TsNode.Content(node.Bytes)) + " "
						}
					}

					def.Signature = sig
				} else if isPassThroughParentNode(node) {
					if verboseLogging {
						fmt.Println("pass through parent node", node.Type)
					}

					start := tsNode.StartByte()

					firstChild := firstDefinitionChild(node)
					if firstChild != nil {
						if verboseLogging {
							fmt.Println("firstChild", firstChild.Type)
						}
						end := firstChild.TsNode.StartByte()
						def.Signature = string(node.Bytes[start:end])

						if verboseLogging {
							fmt.Println("got pass through parent signature", def.Signature)
							fmt.Println("recursing into first child", firstChild.Type)
						}

						def.Children = mapTraditional(node, nil)
					} else {
						if verboseLogging {
							fmt.Println("no first child found", node.Type)
						}
					}

				} else {
					if verboseLogging {
						fmt.Println("not assignment node", node.Type)
						fmt.Println("looking for implementation boundary")
					}
					// Get signature (up to body)
					if body := findImplementationBoundary(node); body != nil {
						if verboseLogging {
							fmt.Println("found implementation boundary", body.Type)
						}

						start := tsNode.StartByte()
						var end uint32
						if tsNode == body.TsNode {
							if verboseLogging {
								fmt.Println("node == body")
							}
							firstChild := firstDefinitionChild(*body)
							if firstChild != nil {
								if verboseLogging {
									fmt.Println("firstChild != nil")
									fmt.Println("firstChild", firstChild.Type)
								}
								end = firstChild.TsNode.StartByte()
							} else {
								if verboseLogging {
									fmt.Println("firstChild == nil")
								}
								end = body.TsNode.EndByte()
							}
						} else {
							end = body.TsNode.StartByte()
						}
						if verboseLogging {
							fmt.Println("start", start)
							fmt.Println("end", end)
						}
						def.Signature = string(node.Bytes[start:end])
						if verboseLogging {
							fmt.Println("got signature", def.Signature)
						}

						// If this is a parent type node, recurse into the body
						if isParentNode(node) {
							if verboseLogging {
								fmt.Println("isParentNode, recursing into body", node.Type)
							}
							def.Children = mapTraditional(*body, &node)
						}
					} else {
						if verboseLogging {
							fmt.Println("no implementation boundary found", node.Type)
						}
						def.Signature = string(node.TsNode.Content(node.Bytes))
					}
				}

				// Get preceding comments
				// no comments for now to minimize tokens
				// def.Comments = getPrecedingComments(node)

				defs = append(defs, def)
			} else {
				if verboseLogging {
					fmt.Println("not definition node", node.Type)
				}
			}

			if !cursor.GoToNextSibling() {
				break
			}
		}
	}

	return defs
}

// // Get preceding comments
// func getPrecedingComments(node Node) []string {
// 	var comments []string
// 	const maxCommentLength = 1000

// 	prevNode := node.TsNode.PrevSibling()
// 	for prevNode != nil {
// 		if !strings.Contains(prevNode.Type(), "comment") {
// 			break
// 		}
// 		comment := string(prevNode.Content(node.Bytes))
// 		if len(comment) > maxCommentLength {
// 			comment = comment[:maxCommentLength] + "..."
// 		}
// 		comments = append([]string{comment}, comments...)
// 		prevNode = prevNode.PrevSibling()
// 	}
// 	return comments
// }

// func mapConfig(node *tree_sitter.Node, content []byte) []Definition {
// 	cursor := tree_sitter.NewTreeCursor(node)
// 	defer cursor.Close()

// 	var walkConfig func(*tree_sitter.Node) []Definition
// 	walkConfig = func(node *tree_sitter.Node) []Definition {
// 		var defs []Definition

// 		// Handle key-value pairs
// 		switch node.Type() {
// 		case "block_mapping_pair": // YAML
// 			if key := node.ChildByFieldName("key"); key != nil {
// 				def := Definition{
// 					Type:      "key",
// 					Signature: string(key.Content(content)),
// 					Line:      int(key.StartPoint().Row) + 1,
// 				}

// 				// Handle nested structures
// 				if val := node.ChildByFieldName("value"); val != nil {
// 					switch val.Type() {
// 					case "block_mapping": // nested YAML map
// 						def.Children = walkConfig(val)
// 					case "block_sequence": // YAML array
// 						// Could track sequences if needed
// 					}
// 				}

// 				defs = append(defs, def)
// 			}
// 		case "pair": // TOML/JSON
// 			// Similar pattern for TOML/JSON
// 		case "field": // CUE/HCL
// 			// Similar pattern for CUE/HCL
// 		}

// 		return defs
// 	}

// 	return walkConfig(node)
// }

func (m *FileMap) String() string {
	var b strings.Builder

	var writeDefinition func(def *Definition, depth int)
	writeDefinition = func(def *Definition, depth int) {
		if def.Type == "svelte-style" {
			b.WriteString("\n")
		}

		// Indent
		if depth > 0 {
			b.WriteString(strings.Repeat("  ", depth))
			b.WriteString("- ")
		}

		// Write signature (for tags, include attrs)
		if def.Type == "tag" {
			// Extract tag name from signature (it's the first word)
			tagName := strings.Fields(def.Signature)[0]
			// Build full representation with attrs
			if len(def.TagAttrs) > 0 {
				if def.TagReps > 1 {
					b.WriteString(fmt.Sprintf("[%dx]", def.TagReps))
				}
				b.WriteString(fmt.Sprintf("%s%s", tagName, strings.Join(def.TagAttrs, "")))
			} else {
				b.WriteString(tagName)
			}
		} else {
			b.WriteString(strings.TrimSpace(def.Signature))
		}
		b.WriteString("\n")

		// Write children with increased depth
		for _, child := range def.Children {
			writeDefinition(&child, depth+1)
		}

		if def.Type == "svelte-script" {
			b.WriteString("\n")
		}
	}

	// Write all top-level definitions
	for _, def := range m.Definitions {
		writeDefinition(&def, 0)
	}

	return b.String()
}