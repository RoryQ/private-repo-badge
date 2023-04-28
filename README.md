# Private Repo Badge

Demo for automatically updating badges that will work in private repos. 

The github action will 
1. Look for the column `Latest Tag` then find the latest semver tag that matches the decoded filename.
2. For each tag a png badge will be generated then uploaded to the configured github release. 


| Package                              | Latest Tag                                                                                                               |
|--------------------------------------|--------------------------------------------------------------------------------------------------------------------------|
| [json2yaml](./json2yaml)             | <img src="https://github.com/RoryQ/private-repo-badge/releases/download/readmebadges/json2yaml.2x.png"  height="20px" /> |
| [yaml2json](./yaml2json)             | <img src="https://github.com/RoryQ/private-repo-badge/releases/download/readmebadges/yaml2json.png" />                   |
| [tools/json2yaml](./tools/json2yaml) | <img src="https://github.com/RoryQ/private-repo-badge/releases/download/readmebadges/tools__json2yaml.png" />            |
| [tools/yaml2json](./tools/yaml2json) | <img src="https://github.com/RoryQ/private-repo-badge/releases/download/readmebadges/tools__yaml2json.png" />            |
