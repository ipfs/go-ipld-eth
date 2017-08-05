# CHANGELOG

## `0.0.3`
  * Improve in-code docs for the plugin loader.
  * Add a vetting script.
  * Work on the input:
    * `eth-block`
      * Accept input in RLP for both block headers or block bodies. (as `raw`).
      * Accept input in JSON.
      * Display common hashes as well as `cid`s on JSON Marshaling.
  * Work on the output:
    * When queried, `eth-block` always responds a block header in RLP.
    * Add resolving (or add tests) for the following elements:
      * `bloom`
      * `coinbase`
      * `parent`
  * Test coverage is `52.7%`.

## `0.0.2`

* First Implementation of vetting tools:
  * `unused`
  * `gosimple`
  * `staticcheck`
  * `golint`
