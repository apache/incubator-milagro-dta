# Milagro-Custody-DTA
---
[![Maintenance](https://img.shields.io/badge/Maintained%3F-yes-green.svg)](https://GitHub.com/Naereen/StrapDown.js/graphs/commit-activity)

Milagro Custody D-TA is a colaborative key management server 

Milagro D-TA facilitates secure and auditable communication between people who to use key pairs (Beneficiaries) and service providers who can keep the secret keys safe (Custodians). It is written in Go and uses REST services based on the GoKit microservices framework: https://gokit.io/, it uses IPFS to create a shared immutable log of transactions and relies on Milagro-Crypto-C for it's crypto.

## Plugins
Milagro D-TA provides a basic set of services for creating identities for actors in the system, and passing encrypted communication between them but it assumes that different service providers will have their own "special sauce" for securely storing public keys, so the vanilla services can be extended using a plugin framework. Two basic plugins are included in this release to give you an idea of how this can be done.
1. *BitcoinPlugin* Will generate a Bitcoin address and reveal the corresponding secret key
2. *SafeGuardSecret* Enables the principal to use a public key to encrypt a string and decrypt it again

## Instalation
To see Milagro D-TA in action clone this repo and run the build script. This has currently been tested on Ubuntu (we'll make a mac version shortly)

```
git clone https://github.com/apache/incubator-milagro-dta.git

cd incubator-milagro-dta

./build_static.sh
```

Alternatively you can run Milagro D-TA in a docker container 

```
git clone https://github.com/apache/incubator-milagro-dta.git

cd incubator-milagro-dta

docker build -t mydta .

```

## Running Test

```
./cmd/servicetester/fulltest.sh
```

## Documentation

You can find documentation for Milagro D-TA in the main [Milagro docs site](https://milagro.apache.org/) 

Which includes a quick start guide that will show you how to get Milagro D-TA to "do stuff"


## Contributing

 Key pairs are becoming key to our online lives, and keeping secret keys safe is a growing industry, we hope to create an ecosystem of custodial service providers who colabroate to make the Internet a safer place for everyone. We are keen to get contributions and feedback from anyone in this space. This is a brand new project so our development processes are still being figured out, but if you have a suggestion, question or wish to make a contribution go ahead raise an issue and someone on the team will get right on it!

## License

Apache License 2.0]

## Credits

Design and Concept... [Brian Spector](https://github.com/spector-in-london)
Core Algorithm and Services... [Chris Morris](https://github.com/fluidjax)
Framework and Refactoring... [Stanislav Mihaylov](https://github.com/smihaylov)
Crypto Genius... [Kealan McCusker](https://github.com/kealan)
Keeper of "The Apache Way"... [John McCane-Whitney](https://github.com/johnmcw)
Prototype and Cat Herding... [Howard Kitto](https://github.com/howardkitto)


















