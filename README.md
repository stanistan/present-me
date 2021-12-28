[(pr)esent-me][1] is an experiment to try to give the author of
a Pull Request a better way to convey why a changeset looks the
way that it does, and how the folks reading and reviewing it
should approach it. ([read more][2])

### How it works

`present-me` uses a [PR review][3]'s comments (and their respective diff)
to create a single "post", or "slides."

---

These are all valid URLs to query for:

- `https://github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746` 

   **Fully qualified Pull Request Review URL (the permalink from GitHub)**
   
- `github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746`
   
   **Dropping the Protocol (https is implicit)**
   
- `stanistan/invoice-proxy/pull/3#pullrequestreview-625362746`
   
   **Dropping the Domain (https://github.com implicit)**
   
- `stanistan/invoice-proxy/pull/3`

   **Dropping the URL Fragment.** This will attempt to find the first PR review
   by the author, and display that if possible.

### Making a Review

1. Create a PR as regular
2. Start a review of your own PR
3. Leave comments on your PR as part of the review.
   - _Comments prefixed with a number will be ordered that way_.
   - Submit your review: <https://github.com/stanistan/invoice-proxy/pull/3#pullrequestreview-625362746>
4. Go to [the website][1], you can put in the URL in the form above (including pullrequestreview),
   or if you put in the PR url (no fragment), it'll try to find the first review
   written by the author of the PR.
   - [generated post](https://present-me.stanistan.dev/stanistan/invoice-proxy/pull/3/625362746/post)
   - [generated slides](https://present-me.stanistan.dev/stanistan/invoice-proxy/pull/3/625362746/slides) using <https://revealjs.com>
   - [generated md](https://present-me.stanistan.dev/stanistan/invoice-proxy/pull/3/625362746/md)

[1]: https://present-me.stanistan.dev
[2]: https://www.stanistan.com/writes/2021/04/13/present-me/
[3]: https://docs.github.com/en/rest/reference/pulls#get-a-review-for-a-pull-request
