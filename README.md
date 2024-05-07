# BlindBitd (wip)

Receive and send functionality. This is the daemon for the BlindBit Wallet. The daemon can be controlled
with [blindbit-cli](./cli/README.md). Still in early testing, only use with funds you can
afford to lose. When started and unlocked the daemon will run continuously in the background using minimal resources
while not processing. When actively scanning a block more resources will be needed. The daemon periodically checks for
new blocks and also listens to Electrum's `blockchain.headers.subscribe`.

IMPORTANT: The wallet data and keys are only encrypted when on disk. For scanning purposes the private keys are
currently always kept in memory. In future there will be a separate spending password that encrypts the spend secret key
and the mnemonic separately.

## Setup

### Requirements

- go 1.21.9 installed
- access to an electrum server
- access to a blindbit style indexing server like [BlindBit Oracle](https://github.com/setavenger/blindbit-oracle)
    - I'm hosting a signet BlindBit-Oracle server here: signet.blindbit.snblago.com:8000
- access to an electrum server
    - I'm hosting a signet Electrum server here: signet.electrum.snblago.com:50001

### Build

There is a makefile from which you can easily build the daemon and the accompanying cli.
The resulting binaries will end up in the project root in a directory called `bin`.
Build both with:

```console
$ make build
```

Build only blindbitd with:

```console
$ make build-daemon
```

Build only blindbit-cli with:

```console
$ make build-cli
```

### Run

The daemon requires a config toml file `blindbit.toml` to be present in its datadir. The socket to the daemon is created
in `<datadir/run/blindbit.socket>`. The path to the socket has to be passed to `blindbit-cli` in order to access the
daemon. The default path for blindbitd is `~/.blindbitd`. For both programs the default path forthe socket is set
to `~/.blindbitd/run/blindbit.socket`.

You can then run with:

```console
$ bin/blindbitd
```

```console
$ bin/blindbitd-cli status
```

## Controlling the daemon

Currently, the daemon is only exposed via a unix socket. The [blindbit-cli](./cli/README.md) controls the flow of the
daemon. On initial startup you can use `createwallet` command to either initialise a new wallet or recover from a
mnemonic. `listaddresses` shows your address. You can use `createtransaction` to send to an address.


## Todo

### Priority 1

- [ ] Bring test coverage to a meaningful level
- [x] Create a coin selector that incorporates the fees
- [x] Binary encoding to reduce bandwidth save time on decoding (protoBuffs)
- [ ] Add Transaction History
- [x] Mark UTXOs as spent (or similar) if used for a transaction
- [ ] Sometimes unlock does not work on first try, needs a restart of the daemon
- [ ] Add gRPC credentials
- [x] Change naming convention of log files
    - It should be easy to determine such that a user does not always have to check the current name

### Priority 2

- [ ] Be able to set Debug level
- [ ] More tests for coin selector
    - Selector seems very accurate, but should rather do +1 sat to exceed fee and don't go below
- [ ] Coin selector allow float fees
- [ ] UTXO export - similar to a backup to avoid rescanning from birthHeight
- [ ] Separate spending password
- [ ] Out-of-band notifications
    - share tweak and tx data directly with the receiver to reduce scanning efforts (follow blindbit standard set for
      the mobile app)
- [ ] Balance checks for UTXOs: account for more than one UTXO per script
- [ ] Expand logging especially on errors
- [ ] Check which panics to keep
- [ ] Automatically make annotation in tx-history if sent to sp-address, not possible to reconstruct in hindsight
- [ ] Don't always add change in coin selector (see todo)

## IPC

- [x] Create New label
    - Give comment
    - Returns label address
- [x] Create Tx and broadcast
- [x] Broadcast raw Tx
- [ ] Pause/Resume scanning
- [x] List UTXOs by label
