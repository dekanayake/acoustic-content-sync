### Acoustic content sync CLI tool.
This tool will help to author contents in a CSV to Acoustic Content headless CMS.
When author has large set of contents to author they can collate those contents in a CSV.
The tool will help them to author the contents in CSV to Acoustic Content.

#### Config yaml
In order to author the content the author need to prepare a config file which maps the 
data field (csv column header) with the Acoustic Content field.

Each element of the config is documented below

##### contentType
The config will start with contentType , this will include all config for a content type
A config can have more than one config type. 

``` yaml
contentType:
  - type: "4c8b4730-7503-485a-9c8e-23af27c61307"
    csvRecordKey : productCode
    name: [productCode,productTitle]
    tags: [moodboard,non-reece,prod,Tile Cloud,2020,November]
    fieldMapping:
     - ...
```


| Config element name        | Value           | Mandatory |
| ------------- |-------------| :-------------:|
|type     | Content type ID  |  Yes |
| csvRecordKey      | The column name from CSV which will use as the key of the record .This will help to find the issues in content when it fails    |  Yes | 
| name | The column names from CSV which will use to generate the name of the content    |  Yes |
| tags | The tags which will added to the content    |  Yes |
| fieldMapping | Mapping configuration of each csv column to Content type field    |  Yes |


###### contentType

A field in Content model can have different data type . Ex: text, number, category, etc..
Based on the data type of the field different confiuration will be applied.

#### Text






