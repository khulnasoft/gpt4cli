"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[777],{5670:(t,e,n)=>{n.r(e),n.d(e,{assets:()=>a,contentTitle:()=>r,default:()=>u,frontMatter:()=>i,metadata:()=>o,toc:()=>c});var l=n(4848),s=n(8453);const i={sidebar_position:2,sidebar_label:"Install"},r="Install Gpt4cli",o={id:"install",title:"Install Gpt4cli",description:"Quick Install",source:"@site/docs/install.md",sourceDirName:".",slug:"/install",permalink:"/install",draft:!1,unlisted:!1,editUrl:"https://github.com/khulnasoft/gpt4cli/tree/main/docs/docs/install.md",tags:[],version:"current",sidebarPosition:2,frontMatter:{sidebar_position:2,sidebar_label:"Install"},sidebar:"tutorialSidebar",previous:{title:"Intro",permalink:"/intro"},next:{title:"Quickstart",permalink:"/quick-start"}},a={},c=[{value:"Quick Install",id:"quick-install",level:2},{value:"Manual install",id:"manual-install",level:2},{value:"Build from source",id:"build-from-source",level:2},{value:"Windows",id:"windows",level:2}];function d(t){const e={a:"a",code:"code",h1:"h1",h2:"h2",p:"p",pre:"pre",...(0,s.R)(),...t.components};return(0,l.jsxs)(l.Fragment,{children:[(0,l.jsx)(e.h1,{id:"install-gpt4cli",children:"Install Gpt4cli"}),"\n",(0,l.jsx)(e.h2,{id:"quick-install",children:"Quick Install"}),"\n",(0,l.jsx)(e.pre,{children:(0,l.jsx)(e.code,{className:"language-bash",children:"curl -sL https://raw.githubusercontent.com/khulnasoft/gpt4cli/main/app/cli/install.sh | bash\n"})}),"\n",(0,l.jsx)(e.h2,{id:"manual-install",children:"Manual install"}),"\n",(0,l.jsxs)(e.p,{children:["Grab the appropriate binary for your platform from the latest ",(0,l.jsx)(e.a,{href:"https://github.com/khulnasoft/gpt4cli/releases",children:"release"})," and put it somewhere in your ",(0,l.jsx)(e.code,{children:"PATH"}),"."]}),"\n",(0,l.jsx)(e.h2,{id:"build-from-source",children:"Build from source"}),"\n",(0,l.jsx)(e.pre,{children:(0,l.jsx)(e.code,{className:"language-bash",children:'git clone https://github.com/khulnasoft/gpt4cli.git\ngit clone https://github.com/khulnasoft-lab/survey.git\ncd gpt4cli/app/cli\ngo build -ldflags "-X gpt4cli/version.Version=$(cat version.txt)"\nmv gpt4cli /usr/local/bin # adapt as needed for your system\n'})}),"\n",(0,l.jsx)(e.h2,{id:"windows",children:"Windows"}),"\n",(0,l.jsxs)(e.p,{children:["Windows is supported via ",(0,l.jsx)(e.a,{href:"https://learn.microsoft.com/en-us/windows/wsl/about",children:"WSL"}),"."]}),"\n",(0,l.jsx)(e.p,{children:"Gpt4cli only works correctly in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell."})]})}function u(t={}){const{wrapper:e}={...(0,s.R)(),...t.components};return e?(0,l.jsx)(e,{...t,children:(0,l.jsx)(d,{...t})}):d(t)}},8453:(t,e,n)=>{n.d(e,{R:()=>r,x:()=>o});var l=n(6540);const s={},i=l.createContext(s);function r(t){const e=l.useContext(i);return l.useMemo((function(){return"function"==typeof t?t(e):{...e,...t}}),[e,t])}function o(t){let e;return e=t.disableParentContext?"function"==typeof t.components?t.components(s):t.components||s:r(t.components),l.createElement(i.Provider,{value:e},t.children)}}}]);