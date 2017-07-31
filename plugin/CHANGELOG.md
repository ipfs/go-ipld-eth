# CHANGELOG

## `0.0.3`
  * Improve in-code docs for the plugin loader.
  * Add a vetting script.
  * Work on the input:
    * `eth-block`
      * Accept input in RLP for both block headers or block bodies. (as `raw`).
      * Accept input in JSON.
  * Work on the output:
    * When queried, `eth-block` always responds a block header in RLP.

## `0.0.2`

* First Implementation of vetting tools:
  * `unused`
  * `gosimple`
  * `staticcheck`
  * `golint`
