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

#### Common field property types
| Config element name | Value |
| ------------- |-------------|
| csvProperty   | The column name in CSV file | 
| acousticProperty  | The mapped acoustic content type field name |
| propertyType | Property type of the content type field |
| staticValue   | Static value , if this property is set csvProperty should removed

#### Supported content type field types
| Property type | 
|---------------|
| number        |
| text          |
|multi-text     |
|formatted-text  |
|toggle         |
|link           |
|date           |
|category       |
|category-part  |
|file           |
|video          |
|image          |
|group          |
|multi-group    |
|reference      |
|multi-reference  |

#### number
Ex:
``` yaml
      - csvProperty: ProductNumber
        acousticProperty: productCode
        propertyType: number
```

#### text
Ex:
``` yaml
      - csvProperty: Product Name
        acousticProperty: productTitle
        propertyType: text
```

#### multi-text
Ex:
``` yaml
      - csvProperty: words
        acousticProperty: word
        propertyType: multi-text
```
#### formatted-text
Ex:
``` yaml
      - csvProperty: description
        acousticProperty: description
        propertyType: formatted-text
```
#### toggle
Ex:
``` yaml
      - acousticProperty: disable
        staticValue: false
        propertyType: toggle
```
#### link
Ex:
``` yaml
      - csvProperty: url
        acousticProperty: url
        propertyType: link
```
#### date
Ex:
``` yaml
      - csvProperty: fromDate
        acousticProperty: fromDate
        propertyType: date
```
#### category
Ex:
``` yaml
      - csvProperty: Category
        acousticProperty: category
        propertyType: category
        categoryName: "Moodboard Non Reece Categories"
```
| Config element name | Value |
|---------------|---------|
| categoryName  | Mapped category name in Acoustic Content |

#### category-part
Ex:
``` yaml
      - csvProperty: categories
        acousticProperty: category
        propertyType: category-part
        categoryName: "Reece AU website category"
        linkToParents: false
```
| Config element name | Value |
|---------------|---------|
| linkToParents  | If true will create category of parents  |

#### file
Ex:
``` yaml
      - csvProperty: MaterialGLB
        acousticProperty: glb
        propertyType: file
        assetName:
          - refCSVProperty: MaterialName
            propertyName: name
        acousticAssetBasePath: "/dxdam/3dplanner/materials"
        assetLocation: "/Users/ekanad/Downloads/Materials"
```
| Config element name | Value |
|---------------|---------|
| assetName.refCSVProperty  | The column name in CSV to generate the asset name |
| assetName.propertyName  | Property name  |
| acousticAssetBasePath  | The base path need to set in Acoustic asset  |
| assetLocation  | The local folder of the assets  |

#### video
``` yaml
      - csvProperty: video
        acousticProperty: video
        propertyType: video
        assetName:
          - refCSVProperty: productId
            propertyName: name
        acousticAssetBasePath: "/dxdam/video"
        assetLocation: "/Users/ekanad/Downloads/videos"
```
| Config element name | Value |
|---------------|---------|
| assetName.refCSVProperty  | The column name in CSV to generate the asset name |
| assetName.propertyName  | Property name  |
| acousticAssetBasePath  | The base path need to set in Acoustic asset  |
| assetLocation  | The local folder of the assets  |

#### image
``` yaml
      - csvProperty: Image
        acousticProperty: image
        propertyType: image
        assetName:
          - refCSVProperty: Product Name
            propertyName: productCode
        profiles: ["ae34cc92-8144-4d78-9660-c7d20abc0817"]
        enforceImageDimension: true
        imageWidth: 1200
        imageHeight: 900
        assetLocation: "/Volumes/dfsroot/Reece/Marketing/9. Digital Marketing/Projects/Bathrooms/2020/Spring/Mood board/Final Images Assets/Swatches/ADP Caesarstone Swatches"
        acousticAssetBasePath: "/dxdam/moodboard/nonreece/prod"
```
| Config element name | Value |
|---------------|---------|
| assetName.refCSVProperty  | The column name in CSV to generate the asset name |
| assetName.propertyName  | Property name  |
| acousticAssetBasePath  | The base path need to set in Acoustic asset  |
| assetLocation  | The local folder of the assets  |
| profiles  | Mapped image profiles |
| enforceImageDimension  | if true image dimension will be update  |
| imageWidth  | Image width to update the dimension.Effective only when enforceImageDimension=true  |
| imageHeight  | Image height to update the dimension.Effective only when enforceImageDimension=true  |

#### group
``` yaml
      - acousticProperty: color
        propertyType: group
        type: "c74414e2-43fc-434b-807e-0ae250478f4f"
        fieldMapping:
          - csvProperty: Color
            acousticProperty: colorpicker
            propertyType: text
          - ....  
```
| Config element name | Value |
|---------------|---------|
| type  | group id in Acoustic content |
| fieldMapping  | field mapping list of the group  |

#### multi-group
``` yaml
      - acousticProperty: color
        propertyType: group
        type: "c74414e2-43fc-434b-807e-0ae250478f4f"
        fieldMapping:
          - csvProperty: Color
            acousticProperty: colorpicker
            propertyType: text
          - ....  
```
| Config element name | Value |
|---------------|---------|
| type  | group id in Acoustic content |
| fieldMapping  | field mapping list of the group  |





