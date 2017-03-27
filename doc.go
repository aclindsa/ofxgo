/*
Package ofxgo seeks to provide a library to make it easier to query and/or
parse financial information with OFX from the comfort of Golang, without having
to deal with marshalling/unmarshalling the SGML or XML. The library does *not*
intend to abstract away all of the details of the OFX specification, which
would be very difficult to do well. Instead, it exposes the OFX SGML/XML
hierarchy as structs which mostly resemble it. For more information on OFX and
to read the specification, see http://ofx.net.

There are three main top-level objects defined in ofxgo. These are Client,
Request, and Response. The Request and Response objects, predictably, contain
representations of OFX requests and responses as structs. Client contains
settings which control how requests and responses are marshalled and
unmarshalled (the OFX version used, client id and version, whether to indent
SGML/XML elements), and provides helper methods for making requests and
optionally parsing the response using those settings.

Every Request object contains a SignonRequest element, called Signon. This
element contains the username, password (or key), and the ORG and FID fields
particular to the financial institution being queried, and an optional ClientUID
field (required by some FIs). Likewise, each Response contains a SignonResponse
object which contains, among other things, the Status of the request. Any status
with a nonzero Code should be inspected for a possible error (using the Severity
and Message fields populated by the server, or the CodeMeaning() and
CodeConditions() functions which return information about a particular code as
specified by the OFX specification).

Each top-level Request or Response object may contain zero or more Messages,
represented by a slice of objects satisfying the Message interface. These
messages are grouped by function into message sets, just as the OFX
specification groups them. Here is a list of the field names of each of these
message sets (each represented by a slices) in the Request/Response objects,
along with the concrete types of Messages they may contain:

Signup:
  AcctInfoRequest/AcctInfoResponse: A listing of the valid accounts for this login

Banking:
  StatementRequest/StatementResponse: The balance (and optionally list of
    transactions) for a bank account

CreditCards:
  CCStatementRequest/CCStatementResponse: The balance (and optionally list of
    transactions) for a credit card

Investments:
  InvStatementRequest/InvStatementResponse: The balance and/or  list of
    transactions for an investment account

Securities:
  SecListRequest/SecListResponse: List securities and their prices, etc.
  SecurityList: The actual list of securities, prices, etc. (even if
    SecListResponse is present, it doesn't contain the security information). Note
    that this is frequently returned with an InvStatementResponse, even if
    SecListRequest wasn't passed to the server.

Profile:
  ProfileRequest/ProfileResponse: Determine the server's capabilities (which
    messages sets it supports, along with individual features)


When constructing a Request, simply append the desired message to the message
set it belongs to. For Responses, it is the user's responsibility to make type
assertions on objects found inside one of these message sets before using them.
*/
package ofxgo
