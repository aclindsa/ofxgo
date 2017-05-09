/*
Package ofxgo seeks to provide a library to make it easier to query and/or
parse financial information with OFX from the comfort of Golang, without having
to deal with marshalling/unmarshalling the SGML or XML. The library does *not*
intend to abstract away all of the details of the OFX specification, which
would be difficult to do well. Instead, it exposes the OFX SGML/XML hierarchy
as structs which mostly resemble it. For more information on OFX and to read
the specification, see http://ofx.net.

There are three main top-level objects defined in ofxgo. These are Client,
Request, and Response. The Request and Response objects represent OFX requests
and responses as Golang structs. Client contains settings which control how
requests and responses are marshalled and unmarshalled (the OFX version used,
client id and version, whether to indent SGML/XML tags, etc.), and provides
helper methods for making requests and optionally parsing the response using
those settings.

Every Request object contains a SignonRequest element, called Signon. This
element contains the username, password (or key), and the ORG and FID fields
particular to the financial institution being queried, and an optional ClientUID
field (required by some FIs). Likewise, each Response contains a SignonResponse
object which contains, among other things, the Status of the request. Any status
with a nonzero Code should be inspected for a possible error (using the Severity
and Message fields populated by the server, or the CodeMeaning() and
CodeConditions() functions which return information about a particular code as
specified by the OFX specification).

Each top-level Request or Response object may contain zero or more messages,
sorted into named slices by message set, just as the OFX specification groups
them. Here are the supported types of Request/Response objects (along with the
name of the slice of Messages they belong to in parentheses):

Requests:
  var r AcctInfoRequest     // (Signup) Request a list of the valid accounts
                            //   for this user
  var r CCStatementRequest  // (CreditCard) Request the balance (and optionally
                            //   list of transactions) for a credit card
  var r StatementRequest    // (Bank) Request the balance (and optionally list
                            //   of transactions) for a bank account
  var r InvStatementRequest // (InvStmt) Request balance, transactions,
                            //   existing positions, and/or open orders for an
                            //   investment account
  var r SecListRequest      // (SecList) Request securities details and prices
  var r ProfileRequest      // (Prof) Request the server's capabilities (which
                            //   messages sets it supports, along with features)

Responses:
  var r AcctInfoResponse     // (Signup) List of the valid accounts for this
                             //   user
  var r CCStatementResponse  // (CreditCard) The balance (and optionally list of
                             //   transactions) for a credit card
  var r StatementResponse    // (Bank) The balance (and optionally list of
                             //   transactions) for a bank account
  var r InvStatementResponse // (InvStmt) The balance, transactions, existing
                             //   positions, and/or open orders for an
                             //   investment account
  var r SecListResponse      // (SecList) Returned as a result of
                             //   SecListRequest, but only contains request
                             //   status
  var r SecurityList         // (SecList) The actual list of securities, prices,
                             //   etc. (sent as a result of SecListRequest or
                             //   InvStatementRequest)
  var r ProfileResponse      // (Prof) Describes the server's capabilities

When constructing a Request, simply append the desired message to the message
set it belongs to. For Responses, it is the user's responsibility to make type
assertions on objects found inside one of these message sets before using them.

For example, the following code would request a bank statement for a checking
account and print the balance:

  import (
    "fmt"
    "github.com/aclindsa/ofxgo"
    "os"
  )

  var client ofxgo.Client // By not initializing them, we accept all default
                          // client values
  var request ofxgo.Request

  // These are all specific to you and your financial institution
  request.URL = "https://ofx.example.com"
  request.Signon.UserID = ofxgo.String("john")
  request.Signon.UserPass = ofxgo.String("hunter2")
  request.Signon.Org = ofxgo.String("MyBank")
  request.Signon.Fid = ofxgo.String("0001")

  uid, err := ofxgo.RandomUID()
  if err != nil {
    fmt.Println("Error creating uid for transaction:", err)
    os.Exit(1)
  }

  statementRequest := ofxgo.StatementRequest{
    TrnUID: *uid,
    BankAcctFrom: ofxgo.BankAcct{
      BankID:   ofxgo.String("123456789"),
      AcctID:   ofxgo.String("11111111111"),
      AcctType: ofxgo.AcctTypeChecking,
    },
  }

  request.Bank = append(request.Bank, &statementRequest)

  response, err := client.Request(request)
  if err != nil {
    fmt.Println("Error requesting account statement:", err)
    os.Exit(1)
  }

  if response.Signon.Status.Code != 0 {
    meaning, _ := response.Signon.Status.CodeMeaning()
    fmt.Printf("Nonzero signon status (%d: %s) with message: %s\n", response.Signon.Status.Code, meaning, response.Signon.Status.Message)
    os.Exit(1)
  }

  if len(response.Bank) < 1 {
    fmt.Println("No banking messages received")
  } else if stmt, ok := response.Bank[0].(*ofxgo.StatementResponse); ok {
    fmt.Printf("Balance: %s %s (as of %s)\n", stmt.BalAmt, stmt.CurDef, stmt.DtAsOf)
  }

More usage examples may be found in the example command-line client provided
with this library, in the cmd/ofx directory of the source.

*/
package ofxgo
