# perfect-gpg-keypair

Makes it really easy to generate a super secure GPG keypair along with a separate signing key to use for the current computer.


## Preface
This was written in MacOS. I have no guarantee that it will work in another OS, but I see no reason why it shouldn't either.


## Why?
Because every time my GPG keypair expires, I find myself having to google how to generate a new one. But mostly because I wanted to play around with golang


## Features
- Interactive (using [gum](https://github.com/charmbracelet/gum))
- Creates a passphrase protected master keypair using RSA4096
- Creates a signing subkey (also using RSA4096) that can be used for signing commits on e.g. GitHub
- Creates a revocation certificate in case of emergency
- Automatically removes the master keypair from your computer (after you have backed them up!) and reimports the signing subkey


## Prerequisites
- Need to have gpg installed


## How to generate a perfect GPG keypair?
Just run the executable and follow the instructions. At one point you will need to take a backup of your keys and the program will halt until you confirm you have done so.


## Resources
- [Creating the perfect gpg keypair](https://alexcabal.com/creating-the-perfect-gpg-keypair)
- [gpg manpages](https://www.gnupg.org/documentation/manpage.html)
- [gum](https://github.com/charmbracelet/gum)

