# rmm-scripts
This repo is used to contains powershell scripts for scripting MS.

### Usage
#### Creating new scripts.
- Run binary `script-generator` with `create` subcommand:
`./script-generator create script-name script-name.ps1 `
- After execution you will get `script-name` directory with `script-name.ps1` and `script-name.json` in it.
- Move `script-name` directory to other scripts (`scripts` directory).
- Create PR with the `script-name` directory added.

#### Updating existing scripts.
- Run binary `script-generator` with `update` subcommand:
`script-generator update <path to changed script-name.ps1> <path to the existing script-name.json>`
- Create PR with the changed `script-name.ps1` and `script-name.json`

### Links
- Build job - http://ci.corp.continuum.net:8080/job/dev_rmm-scripts
- Deploy on DT - http://dtansible.corp.continuum.net:8080/job/dt-ansible-rmm-scripts-deploy/
- Deploy on QA - http://qaansible.corp.continuum.net:8080/job/qa-ansible-rmm-scripts-deploy/
- Build promotion to Stage ready - http://qaansible.corp.continuum.net:8080/job/qa-ansible-promotion-Juno/
- DC-tiket template - https://continuum.atlassian.net/browse/DC-66954
- Checklist example - https://continuum.atlassian.net/wiki/spaces/EN/pages/1177452906/Friday+Deployment+Checklist+23rd+Nov+2018
- DT env. UI - https://control.dtitsupport247.net/qadashb/QuickAccess/NewDesktops
- QA env. UI - https://control.qa.itsupport247.net/qadashb/QuickAccess/NewDesktops
