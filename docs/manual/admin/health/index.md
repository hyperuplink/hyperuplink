# Health

Health is the board's opinion of its own configuration, and it is one fieldset
per problem, each explaining what is wrong and carrying a button that takes you
to the page where it gets fixed.

![The health page](health.webp)

A board with no issues reports that, which is the result you want and the result
you get on a board that was set up carefully.

## Checks

Today it raises one issue, namely an XMPP target configured with
`InsecureSkipVerify`, meaning certificate verification is switched off for that
connection. The flag exists so that a development server holding a self-signed
certificate can be reached at all, however on a live board it means anyone
able to intercept traffic between the board and the XMPP server may present a
certificate of their own and then read or alter everything crossing it, the
credentials included. The button takes you to
[XMPP]({{ manual "admin/comms/xmpp" }}), although the actual fix is in the
board's configuration file rather than on that page, since that is where targets
are defined.

Open the page after every change to the configuration file. The check runs when
the page is opened rather than on a schedule, and it reflects the board as it is
running at that moment.

The page itself: [Administration → Health]({{ hrefTo "admin/health" }})
