"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[813],{770:(e,n,i)=>{i.r(n),i.d(n,{assets:()=>l,contentTitle:()=>t,default:()=>u,frontMatter:()=>r,metadata:()=>a,toc:()=>c});var s=i(4848),o=i(8453);const r={sidebar_position:9,sidebar_label:"Collaboration / Orgs"},t="Collaboration and Orgs",a={id:"core-concepts/orgs",title:"Collaboration and Orgs",description:"While so far Gpt4cli is mainly focused on a single-user experience, we plan to add features for sharing, collaboration, and team management in the future, and some groundwork has already been done. Orgs are the basis for collaboration in Gpt4cli.",source:"@site/docs/core-concepts/orgs.md",sourceDirName:"core-concepts",slug:"/core-concepts/orgs",permalink:"/core-concepts/orgs",draft:!1,unlisted:!1,editUrl:"https://github.com/khulnasoft/gpt4cli/tree/main/docs/docs/core-concepts/orgs.md",tags:[],version:"current",sidebarPosition:9,frontMatter:{sidebar_position:9,sidebar_label:"Collaboration / Orgs"},sidebar:"tutorialSidebar",previous:{title:"Background Tasks",permalink:"/core-concepts/background-tasks"},next:{title:"Providers",permalink:"/models/model-providers"}},l={},c=[{value:"Multiple Users",id:"multiple-users",level:2},{value:"Domain Access",id:"domain-access",level:2},{value:"Invitations",id:"invitations",level:2},{value:"Joining an Org",id:"joining-an-org",level:2},{value:"Listing Users and Invites",id:"listing-users-and-invites",level:2},{value:"Revoking Users and Invites",id:"revoking-users-and-invites",level:2}];function d(e){const n={code:"code",h1:"h1",h2:"h2",p:"p",pre:"pre",strong:"strong",...(0,o.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(n.h1,{id:"collaboration-and-orgs",children:"Collaboration and Orgs"}),"\n",(0,s.jsxs)(n.p,{children:["While so far Gpt4cli is mainly focused on a single-user experience, we plan to add features for sharing, collaboration, and team management in the future, and some groundwork has already been done. ",(0,s.jsx)(n.strong,{children:"Orgs"})," are the basis for collaboration in Gpt4cli."]}),"\n",(0,s.jsx)(n.h2,{id:"multiple-users",children:"Multiple Users"}),"\n",(0,s.jsxs)(n.p,{children:["Orgs are helpful already if you have multiple users using Gpt4cli in the same project. Because Gpt4cli outputs a ",(0,s.jsx)(n.code,{children:".gpt4cli"})," file containing a bit of non-sensitive config data in each directory a plan is created in, you'll have problems with multiple users unless you either get each user into the same org or put ",(0,s.jsx)(n.code,{children:".gpt4cli"})," in your ",(0,s.jsx)(n.code,{children:".gitignore"})," file. Otherwise, each user will overwrite other users' ",(0,s.jsx)(n.code,{children:".gpt4cli"})," files on every push, and no one will be happy."]}),"\n",(0,s.jsx)(n.h2,{id:"domain-access",children:"Domain Access"}),"\n",(0,s.jsx)(n.p,{children:"When starting out with Gpt4cli and creating a new org, you have the option of automatically granting access to anyone with an email address on your domain."}),"\n",(0,s.jsx)(n.h2,{id:"invitations",children:"Invitations"}),"\n",(0,s.jsxs)(n.p,{children:["If you choose not to grant access to your whole domain, or you want to invite someone from outside your email domain, you can use ",(0,s.jsx)(n.code,{children:"gpt4cli invite"}),":"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"gpt4cli invite\n"})}),"\n",(0,s.jsx)(n.h2,{id:"joining-an-org",children:"Joining an Org"}),"\n",(0,s.jsxs)(n.p,{children:["To join an org you've been invited to, use ",(0,s.jsx)(n.code,{children:"gpt4cli sign-in"}),":"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"gpt4cli sign-in\n"})}),"\n",(0,s.jsx)(n.h2,{id:"listing-users-and-invites",children:"Listing Users and Invites"}),"\n",(0,s.jsxs)(n.p,{children:["To list users and pending invites, use ",(0,s.jsx)(n.code,{children:"gpt4cli users"}),":"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"gpt4cli users\n"})}),"\n",(0,s.jsx)(n.h2,{id:"revoking-users-and-invites",children:"Revoking Users and Invites"}),"\n",(0,s.jsxs)(n.p,{children:["To revoke an invite or remove a user, use ",(0,s.jsx)(n.code,{children:"gpt4cli revoke"}),":"]}),"\n",(0,s.jsx)(n.pre,{children:(0,s.jsx)(n.code,{className:"language-bash",children:"gpt4cli revoke\n"})})]})}function u(e={}){const{wrapper:n}={...(0,o.R)(),...e.components};return n?(0,s.jsx)(n,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},8453:(e,n,i)=>{i.d(n,{R:()=>t,x:()=>a});var s=i(6540);const o={},r=s.createContext(o);function t(e){const n=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function a(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:t(e.components),s.createElement(r.Provider,{value:n},e.children)}}}]);