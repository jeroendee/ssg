# Implementation Tasks

## 1. Model Layer
- [x] 1.1 Add `OGImage string` field to `model.Site` struct
- [x] 1.2 Add `OGImage string` field to `model.Config` struct

## 2. Configuration Layer
- [x] 2.1 Add `OGImage string` to `yamlConfig.Site` struct with `yaml:"ogImage"`
- [x] 2.2 Map yamlConfig.Site.OGImage to model.Config.OGImage in Load()

## 3. Renderer Layer
- [x] 3.1 Update `ogImageURL()` to prefer `site.OGImage` over `site.Logo`

## 4. Testing
- [x] 4.1 Add config test for ogImage parsing
- [x] 4.2 Add renderer test for OGImage priority over Logo
- [x] 4.3 Add test for Logo fallback when OGImage empty

## 5. Documentation
- [x] 5.1 Update dev/ssg.yaml with ogImage example
