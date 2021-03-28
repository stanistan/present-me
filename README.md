# present-me

We get our context set up so we can chat with the client, this is
something that's got to be pretty simple, maybe we do this with a server
at some point, but to start with running this as a command that takes
some args would be pretty neato, and it can dump out some markdown.

If this runs as a server we will then want to do take this markdown,
render it, etc.

We can split this out into different `cmd/` s.

The general behavior is this:

1. We get the Review
2. We get the ReviewComments
3. We get the Files

Determine the order of the markdown to print:

0. The PR Title

1. the Review.Body

2. The ReviewComments.Body in order that they appear, the ordering of the
comments is by 1) if they start with a number!, and then by the rest.
this should be fine since we're not going to have _too many_ comments, this
is a thing a person would do on their own. Of course.

Each one will also have an associated file path. Keep track of this so we can know
what the rest of the file paths are, at the end of the PR?

Will probably want to play with the representation of these (with _changes_, and links
to the original source, etc).

3. Stuff about the rest of the files?
