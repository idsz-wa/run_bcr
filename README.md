## Description
This Wrapper is a package designed for analyzing single-cell BCR data from rabbits. It can analyze raw FASTQ data and output BCR clone information for each cell.

### Software Requirements
* [MIXCR Version > 4.5 ](https://github.com/milaboratory/mixcr)
  
#### Usage

* Download 

```
git clone https://github.com/idsz-wa/run_bcr
cd run_bcr
chmod +x run_bcr_linux 
```

* Edit configuration file 
  
```
##define file location
Reads:
 - class: File
   location: "/mnt/d/bcr_test/t1.r1.fq.gz" # fastq R1 location 

 - class: File
   location: "/mnt/d/bcr_test/t1.r2.fq.gz"  # fastq R2 location 

##define tools path
mixcr_path:
  class: File
  location: "/home/idsz/miniconda3/envs/py3/bin/mixcr" # mixcr location, Please ensure that you have obtained a license for MixCR.

##define  imgt_json path:
mixcr_json:
  class: File
  location: "/mnt/d/bcr_test/imgt.202312-3.sv8_rabbit.json" # imgt_json location 
version: 1
threads: 16
```
* Run 
```
run_bcr -c <config_file> -w <working_directory> 
```
