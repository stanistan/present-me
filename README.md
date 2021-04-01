# present-me

**Present Me** is an application should make presenting
a Pull Request easier.

Generally, when we submit PRs they are read top to bottom,
with little control to the _author_ of the PR in how they
want to present the information (outside of the PR description).

`present-me` uses Github's Review Comments as an annotation method
to generate a markdown presentation as well as slides in [remarkjs][1].

**This requires Read Access to work.**

## TODO

- [ ] multiple kinds of auth (currently uses github application private key)
- [ ] different stylesheets
- [ ] list files in the PR (that are not commented)
- [ ] better logging and error handling

[1]: http://remarkjs.com
