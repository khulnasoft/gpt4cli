package prompts

import (
	"fmt"
	"strings"

	"gpt4cli-server/syntax"
	shared "gpt4cli-shared"
)

type ValidationPromptParams struct {
	Path         string
	Original     string
	Desc         string
	Proposed     string
	Diff         string
	Reasons      []syntax.NeedsVerifyReason
	SyntaxErrors []string
}

// GetValidationReplacementsXmlPrompt constructs the complete prompt string for XML responses.
func GetValidationReplacementsXmlPrompt(params ValidationPromptParams) string {
	reasons := params.Reasons
	syntaxErrs := params.SyntaxErrors
	path := params.Path
	original := params.Original
	desc := params.Desc
	proposed := params.Proposed
	diff := params.Diff

	s := getBuildPromptHead(path, original, desc, proposed)

	s += fmt.Sprintf(
		`
Diff of applied changes:
>>>
%s
<<<

`,
		diff,
	)

	var parts []string

	reasonMap := map[syntax.NeedsVerifyReason]string{
		syntax.NeedsVerifyReasonAmbiguousLocation: "Changes were applied to an ambiguous location. This may indicate incorrect anchor spacing/indentation, wrong anchor ordering, or missing context.",
		syntax.NeedsVerifyReasonCodeRemoved:       "Code was removed or replaced. Verify if this was intentional according to the plan.",
		syntax.NeedsVerifyReasonCodeDuplicated:    "Code may have been duplicated. Verify if this was intentional according to the plan.",
	}

	for _, reason := range reasons {
		if msg, ok := reasonMap[reason]; ok {
			parts = append(parts, msg)
		}
	}

	if len(syntaxErrs) > 0 {
		parts = append(parts, fmt.Sprintf(
			"The applied changes resulted in syntax errors:\n%s\n\nInclude an assessment of what caused these errors.",
			strings.Join(syntaxErrs, "\n"),
		))
	}

	s += strings.Join(parts, "\n\n")

	s += `
## Validation

Your first task is to examine whether the changes were applied as described in the proposed changes explanation. Do NOT evaluate:
- Code quality
- Missing imports
- Unused variables
- Best practices
- Potential bugs
- Syntax (unless syntax errors have been previously specified and you are determining the cause of the syntax errors)

Your evaluation should ONLY assess:
a. Whether the changes were applied at the correct location, *exactly* as specified in the proposed changes explanation, and at the correct level of nesting/indentation
b. Whether the changes included *all* the specified additions/modifications
c. Whether *any* unintended changes were made to surrounding code
d. Whether *any* specified code was accidentally removed or duplicated
e. Any syntax errors that have been previously specified

--

Line numbers prefixed with 'g4c-' are included in the original file. Line numbers prefixed with 'g4c-new-' are included in the proposed changes. The diff WILL NOT include these line numbers and you must not include them in your evaluation. You must ignore them completely.

--

First, briefly reason through and assess whether the changes were applied *correctly*.
You MUST include reasoning–do not skip this step.

If the changes were applied *correctly*, you MUST output a <Gpt4cliCorrect/> tag, followed by a <Gpt4cliFinish/> tag, then end your response, like this:

<Gpt4cliCorrect/>
<Gpt4cliFinish/>

--

If the changes were applied *incorrectly*, first assess what went wrong in your reasoning, and briefly strategize on how these issues can be avoided when you generate replacements. You MUST include reasoning–do not skip this step.

Next, you MUST output a <Gpt4cliIncorrect/> tag, and then proceed to output the <Gpt4cliComments/> tag and the <Gpt4cliReplacements/> tag with at least one <Replacement> element (see below for details). Example:

<Gpt4cliIncorrect/>
<Gpt4cliComments>
...
</Gpt4cliComments>
<Gpt4cliReplacements>
  <Replacement>
    <Old>...</Old>
    <New>...</New>
  </Replacement>
</Gpt4cliReplacements>

--

## Comments

Next, if the changes were applied *incorrectly*: 

` + CommentClassifierPrompt + `

--

## Replacements

Next, if the changes were applied *incorrectly*, you must analyze the *original file* and the *proposed updates* and output a <Gpt4cliReplacements> element that applies the changes described in the *proposed updates* to the *original file* in order to produce a final, valid resulting file with all changes correctly applied.

CRITICALLY IMPORTANT: When applying changes with replacements, NO REFERENCE COMMENTS CAN BE PRESENT IN THE RESULTING FILE. All reference comments (as listed in the <Gpt4cliComments> element above) ABSOLUTELY MUST be replaced with the code they refer to in the *original file*.

Now output a <Gpt4cliReplacements> element that contains all the replacements needed to correctly apply the changes described in the *proposed updates* to the *original file*. The <Gpt4cliReplacements> element MUST contain at least one <Replacement> element.

For each replacement, use a <Replacement> element with the following structure:

<Replacement>
  <Old>...</Old>  
  <New>...</New>
</Replacement>

The <Old> element must contain the *exact* original code that will be replaced. *Every* character in the <Old> element must be present in the original file. You MUST include line numbers prefixed with 'g4c-' in the <Old> element (NOT with 'g4c-new-'). Every line in the <Old> element must exactly match a line in the original file, including spacing, indentation, newlines, and the 'g4c-' line number. <Old> MUST NOT contain any partial lines, only complete lines.

The <New> element must contain ALL the new code that will replace the code in <Old>. It must contain complete lines only (no partial lines). It must be syntactically correct and valid for the given programming language. It MUST NOT contain any line numbers. It MUST NOT contain any reference comments listed in the <Gpt4cliComments> element. ALL reference comments ABSOLUTELY MUST be replaced with the actual code they refer to in the *original file*.

Apply changes intelligently *in order* to avoid syntax errors, breaking code, or removing code from the original file that should not be removed. Consider the reason behind the update and make sure the result is consistent with the intention of the plan.

Pay *EXTREMELY close attention* to opening and closing brackets, parentheses, and braces. Never leave them unbalanced when the changes are applied. Also pay *EXTREMELY close attention* to newlines and indentation. Make sure that the indentation of the new code is consistent with the indentation of the original code, and syntactically correct.

Replacements must be ordered according to their position in the file. Each <Old> block must come after the previous block in the file. Replacements MUST NOT overlap. If a replacement is dependent on another replacement or intersects with it, group those replacements together into a single <Gpt4cliReplacement> block.

You ABSOLUTELY MUST NOT overwrite or delete code from the original file unless the plan *clearly intends* for the code to be overwritten or removed. Do NOT replace a full section of code with only new code unless that is the clear intention of the plan. Instead, merge the original code and the proposed updates together intelligently according to the intention of the plan.

--

Example responses:

1. Changes Applied Correctly:

## Evaluate Diff
The new function 'someFunction' was correctly added to the end of the file, with proper indentation and spacing.

<Gpt4cliCorrect/>
<Gpt4cliFinish/>

2. Changes Applied Incorrectly:

## Evaluate Diff
The new function 'someFunction' was incorrectly added to the end of the file - it was inserted with wrong indentation.

<Gpt4cliIncorrect/>

<Gpt4cliComments>
g4c-new-42: // Update the user
Evaluation: Describes the change being made. Not a reference.
Reference: false

g4c-new-44: // ... existing code ...
Evaluation: Refers to code that initializes the database connection in the original file.
Reference: true
</Gpt4cliComments>

<Gpt4cliReplacements>
  <Replacement>
    <Old>
      g4c-42: func someFunction() {
      g4c-43:   connectToDatabase()
      g4c-44: }
    </Old>
    <New>
      func someFunction() {
        err := connectToDatabase()
        if err != nil {
          log.Printf("error: %v", err)
          return
        }
        processData()
      }
    </New>
  </Replacement>
</Gpt4cliReplacements>

IMPORTANT RULES:
1. If your evaluation finds ANY issues, you MUST use <Gpt4cliIncorrect/> followed by a <Gpt4cliComments> element and a <Gpt4cliReplacements> element with at least one <Replacement> element.
2. If your evaluation finds NO issues, you MUST use <Gpt4cliCorrect/> then a <Gpt4cliFinish/> element. Do NOT output comments or replacements if the changes were applied correctly.
3. In replacements, every line in the <Old> element MUST exactly match a line in the original file and MUST begin with the line number with a 'g4c-' prefix (NOT with a 'g4c-new-' prefix).
4. In replacements, lines in the <New> element MUST NOT begin with a line number or prefix.
5. Always include reasoning in a '## Evaluate Diff' section prior to outputting the <Gpt4cliCorrect/> or <Gpt4cliIncorrect/> tags.
`

	return s
}

// getBuildPromptHead describes the original file and proposed changes
func getBuildPromptHead(filePath, preBuildState, desc, proposed string) string {
	preBuildStateWithLineNums := shared.AddLineNums(preBuildState)
	proposedWithLineNums := shared.AddLineNumsWithPrefix(proposed, "g4c-new-")

	return fmt.Sprintf(
		`Path: %s

Original file (with line nums):
>>>
%s
<<<

Proposed changes explanation (with line nums):
>>>
%s
<<<

Proposed changes:
>>>
%s
<<<
`,
		filePath,
		preBuildStateWithLineNums,
		desc,
		proposedWithLineNums,
	)
}
