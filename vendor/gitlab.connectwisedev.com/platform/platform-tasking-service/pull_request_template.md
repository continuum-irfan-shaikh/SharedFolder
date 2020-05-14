## Description

Please add description what your code do or improve

## Checklist

During creation of new PR please make sure that you walked through list of PR-related questions and marked all true statements.

### Code hygiene

- [ ] Does the project compile?
- [ ] Does code pass all linting and sanity checks locally? 
- [ ] Have you removed all unused, dead or debugging code?

### Code organisation and correctness

- [ ] If an algorithm is material to your code changes, have you ensured that its as efficient and correct as possible?
- [ ] If you are adding logic as part of your PR – are you adding it in the right level?

### Documentation

- [ ] Branch name and PR's title correspond to [Git Source Control Process For Developers](https://continuum.atlassian.net/wiki/spaces/EN/pages/95192705/Git+Source+Control+Process+For+Developers#GitSourceControlProcessForDevelopers-Namingconventions)
- [ ] Is User Story or Technical Issue is well documented and include sane acceptance criteria?
- [ ] Have you made comments in the code you are touching that will be necessary to explain what you did?
- [ ] If PR made some changes in application's REST API - does it includes changes for Swagger documentation?

### Testing

- [ ] Have you manually tested the behavior you are affecting in the PR?
- [ ] Have you built tests that cover the code touched by your PR?

### Dependencies

- [ ] If this PR introducing new dependencies, does they added to `glide.yaml`? 
