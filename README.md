# Zuri NFT CLI Script

This is a Golang script that converts the HNG NFT csv file to JSON files
according to the CHIP-0007 standard for CHIA wallet.

## Requires

`golang 1.17` 

## How to Setup

- Clone this repository using:

```shell
git clone https://github.com/Emmrys-Jay/zuri-nft-script.git
```

- Change into the directory using:

```shell
cd zuri-nft-script
```

- Build this project using:
```shell
go build
```
  - A binary file named `zuri-nft-script` will be generated in your root directory.

- Copy this generated binary file to a directory where the csv file is located, 
or copy the CSV file to your root directory.

- Run binary with command line flag in the form:

```shell
.\zuri-nft-script -csv <filename.csv>
```
where filename is the name of the csv file you just copied.

- The script generates a new CSV file `filename.output.csv` with an extra SHA256 field, 
and a folder `nft-jsons` containing the CHIP-0007 json format files for each nft.




