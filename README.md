# Zuri NFT CLI Script

This is a Go (Golang) script that converts the HNG NFT csv file to JSON files
according to the CHIP-0007 standard for CHIA wallet.

## Requires

`golang 1.17` or higher

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

- The script generates a new CSV file `filename.output.csv` with an extra SHA256 field, 
and a folder `nft-jsons` containing the CHIP-0007 json format files for each nft.

## CHIP-0007 JSON data sample

```shell
{
    "format": "CHIP-0007",
    "name": "Pikachu",
    "description": "Electric-type Pokémon with stretchy cheeks",
    "minting_tool": "SuperMinter/2.5.2",
    "sensitive_content": false,
    "series_number": 22,
    "series_total": 1000,
    "attributes": [
        {
            "trait_type": "Species",
            "value": "Mouse"
        },
        {
            "trait_type": "Color",
            "value": "Yellow"
        },
        {
            "trait_type": "Friendship",
            "value": 50,
            "min_value": 0,
            "max_value": 255
        }
    ],
    "collection": {
        "name": "Example Pokémon Collection",
        "id": "e43fcfe6-1d5c-4d6e-82da-5de3aa8b3b57",
        "attributes": [
            {
                "type": "description",
                "value": "Example Pokémon Collection is the best Pokémon collection. Get yours today!"
            },
            {
                "type": "icon",
                "value": "https://examplepokemoncollection.com/image/icon.png"
            },
            {
                "type": "banner",
                "value": "https://examplepokemoncollection.com/image/banner.png"
            },
            {
                "type": "twitter",
                "value": "ExamplePokemonCollection"
            },
            {
                "type": "website",
                "value": "https://examplepokemoncollection.com/"
            }
        ]
    },
    "data": {
        "example_data": "VGhpcyBpcyBhbiBleGFtcGxlIG9mIGRhdGEgdGhhdCB5b3UgbWlnaHQgd2FudCB0byBzdG9yZSBpbiB0aGUgZGF0YSBvYmplY3QuIE5GVCBhdHRyaWJ1dGVzIHdoaWNoIGFyZSBub3QgaHVtYW4gcmVhZGFibGUgc2hvdWxkIGJlIHBsYWNlZCB3aXRoaW4gdGhpcyBvYmplY3QsIGFuZCB0aGUgYXR0cmlidXRlcyBhcnJheSB1c2VkIG9ubHkgZm9yIGluZm9ybWF0aW9uIHdoaWNoIGlzIGludGVuZGVkIHRvIGJlIHJlYWQgYnkgdGhlIHVzZXIu"
    }
}
```


