# Gemini Jobcoin Mixer

## Requirements

- Docker

## Setup and Use

build the docker image
```bash
make build
```

enter the docker image to run commands
```bash
make sh
```

initialize the db
```bash
./gemini_jobcoin_mixer initDb
```

create some pathways (mixings)
```bash
./gemini_jobcoin_mixer new <address_1> <address_2> ...
```

Go to https://jobcoin.gemini.com/carnation and deposit some coins into the deposit address.

run the mixer
```bash
./gemini_jobcoin_mixer poll
```

the mixer is meant to be run multiple times, pruning a bit of the pool into the member accounts each time. This is to randomize the public data to obfuscate reverse engineering the pathways.

you can see available commands by running
```bash
./gemini_jobcoin_mixer help
```

## Discussion

I enjoyed this project. I wish I had more time to give it, but I had to timebox the limited amount of time I had. Because of this, there are things I didn't get to that I very much wish I did, most notably writing unit test.

Other improvements I would like to make:
- Test for various long numbers for coinage.
- Create a server that my CLI runs against that can have an automatic polling on a set schedule.
- DRY the code up, notably with the Table class.
- Better handling of exceptions.
- Optimized pathway lookup based upon "last checked" timestamp. Currently as is, there is O(n) processing of the pathways.
- Polish the Docker and Make files.
- You'll see miscellanious `TODOs` peppered throughout the code. Those I would definitely want to get to, and when at work I don't push `TODOs` with my PRs.

I spent most of my time on the language itself. This is only my second GoLang app, and translating the Domain Driven Design approach to the constraints of GoLang was a challenge. This is most pronounced with the lack of classes in the language, with packages being an almost-replacement. Functional approaches also don't work well with the language I found.

I used dependency injection where I could, and the native `wire` tool as a service provider to clean the code a bit.
