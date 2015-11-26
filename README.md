# Rob DNS Server

A simple DNS server.

Goals:

* DNS forwarding
  * You ask for an address, if it doesn't know it, it looks it up from someone else.
* Caching
  * You ask for an address, we've already looked it up, we give back the answer.
* Authoritative zone
  * We own a bunch of domains. Maybe with a custom TLD. You can add/remove addresses at will.

Basic tenants:

* Concurrency
  * Handle many connections/queries without trying too hard.
* Respecting timeouts
  * It would be nice to respond in a timely fashion, if appropriate.
