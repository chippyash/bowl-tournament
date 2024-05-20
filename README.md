# Bowl Tournament

A cli utility to create bowls fixture pictorial representation e.g.

![docs/demo-graph.puml](docs/demo-graph.png)

## Pre Installation
Install Graphviz

Debian: `sudo apt install graphviz`

Fedora: `sudo dnf install graphviz`

Install Java

Debian: `sudo apt install default-jre`

Fedora: see: https://docs.fedoraproject.org/en-US/quick-docs/installing-java/

## Manual Installation
Download the `bowl-tournament.tar.gz` file from the dist directory of the repo. 

Unzip the distribution package

 - `tar -xf bowl-tournament.tar.gz` or use your file manager

Move the `bowl-tournament` file to a directory in your path, usually `/usr/local/bin`

 - `sudo mv bowl-tournament /usr/local/bin`

Move the `plantuml.jar` file to `/usr/share/bowl-tournament/plantuml.jar`

 - `sudo mkdir -p /usr/share/bowl-tournament`
 - `sudo mv plantuml.jar /usr/share/bowl-tournament/plantuml.jar`

## Install using your package manager
### Debian based distributions
Download the `bowl-tournament*.deb` from the dist directory of the repo

run `sudo dpkg -i bowl-tournament*.deb`

### Fedora based distributions
Download the `bowl-tournament*.rpm` from the dist directory of the repo

run `sudo yum localinstall bowl-tournament*.rpm`

## Check your installation
make sure you can run the program

`bowl-tournament help` should display something similar to

```text
Bowl tournament utility

Usage:
   bowl-tournament <input> <tournament> <game> {flags}
   bowl-tournament <command> {flags}

Commands: 
   help                          displays usage informationn
   version                       displays version number

Arguments: 
   input                         input csv file
   tournament                    name of tournament
   game                          name of game

Flags: 
   -g, --groupdates              group dates in the round data (default: false)
   -h, --help                    displays usage information of the application or a command (default: false)
   -o, --output                  output png file (default: false)
   -t, --theme                   the Plantuml theme to use. See https://plantuml.com/theme (default: _none_)
   -v, --version                 displays version number (default: false)
```

Create a working directory where you will put all your csv files for processing. 

## Usage

### CSV Files
The utility requires csv files with the following headings:

```text
round,match,play_by,play_on,time,home_participant,away_participant
```
The actual heading names don't matter as they are not used, but the heading line must exist and the content of the data
lines must conform to their intended use.

For example, this is for a singles competition that has play by dates
```text
round,match,play_by,play_on,time,home_participant,away_participant
1,K1R1M2,02/06/24,,,Player 1,Player 2
1,K1R1M4,02/06/24,,,Player 3,Player 4
1,K1R1M6,02/06/24,,,Player 5,Player 6
1,K1R1M8,02/06/24,,,Player 7,Player 8
1,K1R1M10,02/06/24,,,Player 9,Player 10
1,K1R1M12,02/06/24,,,Player 11,Player 12
1,K1R1M14,02/06/24,,,Player 13,Player 14
1,K1R1M16,02/06/24,,,Player 15,Player 16
2,K1R2M1,23/06/24,,,Player 17,Winner of K1R1M2
2,K1R2M2,23/06/24,,,Player 18,Winner of K1R1M4
2,K1R2M3,23/06/24,,,Player 19,Winner of K1R1M6
2,K1R2M4,23/06/24,,,Player 21,Winner of K1R1M8
2,K1R2M5,23/06/24,,,Player 22,Winner of K1R1M10
2,K1R2M6,23/06/24,,,Player 23,Winner of K1R1M12
2,K1R2M7,23/06/24,,,Player 24,Winner of K1R1M14
2,K1R2M8,23/06/24,,,Player 25,Winner of K1R1M16
3,K1R3M1,21/07/24,,,Winner of K1R2M1,Winner of K1R2M2
3,K1R3M2,21/07/24,,,Winner of K1R2M3,Winner of K1R2M4
3,K1R3M3,21/07/24,,,Winner of K1R2M5,Winner of K1R2M6
3,K1R3M4,21/07/24,,,Winner of K1R2M7,Winner of K1R2M8
4,K1R4M1,18/08/24,,,Winner of K1R3M1,Winner of K1R3M2
4,K1R4M2,18/08/24,,,Winner of K1R3M3,Winner of K1R3M4
5,K1R5M1,,,,Winner of K1R4M1,Winner of K1R4M2

```
The match designator can be anything, but you should look for consistency.

Dates **must** be in the format dd/mm/yy

Where you indicate that one of the participants is the result of a match in a previous round, it **must be written** as 'Winner of match'

And this for a triples competition with Play on dates and times
```text
round,match,play_by,play_on,time,home_participant,away_participant
1,K1R1M1,,11/07/24,6.30pm,"Team A","Team B"
1,K1R1M2,,11/07/24,6.30pm,"Team C","Team  D"
1,K1R1M3,,11/07/24,6.30pm,"Team E","Team F"
1,K1R1M4,,11/07/24,6.30pm,"Team G","Team H"
2,K1R2M1,,01/08/24,6.30pm,Winner of K1R1M1,Winner of K1R1M2
2,K1R2M2,,01/08/24,6.30pm,Winner of K1R1M3,Winner of K1R1M4
3,K1R3M1,,,,Winner of K1R2M1,Winner of K1R2M2
```
**NB** The participant names are surrounded by speach quotes as they can include commas.

### Creating Image Files

**NB** *To regenerate the image files, you must first delete any image (.png) files that you require regenerating. The program 
will not overwrite existing image files.*

At the terminal, navigate to your csv directory.

For each of your csv files run the command, e.g. for file MensTrips.csv in the FooBar Tournament:

`bowl-tournament "MensTrips.csv" "FooBar Tournament" "Mens Triples" -g -o`

This will create a `MensTrips.puml` file, which is a [Plantuml](https://plantuml.com/class-diagram) class diagram definition
and a `MensTrips.png` file which is the image file you can insert into your game notice page.

You can view the puml file by typing 

`java -jar /usr/share/bowl-tournament/plantuml.jar -gui` 

in your terminal in your working directory.

If you don't use the '-o' flag, run

`java -jar /usr/share/bowl-tournament/plantuml.jar ./*.puml` 

This will create png image files for all the plantuml definitions. 
Insert the images into a LibreOffice document and print. 

### Process an entire directory

**Important** You need to name your csv files in a consistent manner and include the game name in the file name e.g.

 - Outdoor Championship 2024 Ladies Australian Pairs.csv
 - Outdoor Championship 2024 Mens Australian Pairs.csv
 - Outdoor Championship 2024 Mixed Singles.csv

Create a regex pattern to pull out the game name from the file name. e.g. `.*2024 (.*)` will work on the above file names.
Note the group capture `(.*)`. Your regex must contain one group capture to find the game name.
You can use [https://regex101.com/](https://regex101.com/) to test your pattern.

Run the bowl-tournament command replacing the file name with the directory you want to process and the game with your regex pattern.
Assuming you are currently in your terminal in your working directory then the following will work:

`bowl-tournament $PWD "RTBC Outdoor Championship 2024" ".*2024 (.*)" -g -o -l`

## For development
**Read the [CONTRIBUTING](CONTRIBUTING.md) guidelines**

The utility is written in Go V1.21+
 - fork this repo and then change directory to a clone of your fork.
 - create a new bug or feature branch
 - install the `make` tool if not already done.
 - run `make help` for list of available make commands.
 - run `make deploy-local` to build and deploy the utility locally.
 - if you want to build the rpm and deb packages:
   - run `make getfpm` to install FPM and its dependencies. This needs doing only once.
   - run `make package` to create the deb and rpm packages.

The source is in bowl-tournament.go.

Before submitting a pull request run `make package`.

Create a PR from your fork to the master branch of this repository

## Help Required
If you know how to build and deploy this application for Windows and/or Macintosh, please consider making a pull request.

## License
This software is licensed under the MIT License. See LICENSE.