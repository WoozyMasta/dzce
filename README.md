# dzce

Go package with typed models and codecs for DayZ Central Economy (CE) files.

## Install

```bash
go get github.com/woozymasta/dzce@latest
```

## What It Covers

Main CE families:

* XML (`types`, `events`, `economy`, `globals`, `messages`, and others)
* JSON (`cfggameplay.json`, `cfgeffectarea.json`, etc.)
* binary `areaflags.map`
* CE project `zg-config` (`mapname.xml`)

Detailed areaflags/zg-config format notes in [AREAFLAGS.md](./AREAFLAGS.md).

## Quick Start

Kind-based decode/encode:

```go
kind := dzce.DetectKind("db/types.xml")

value, err := dzce.Decode(kind, input)
if err != nil {
    return err
}

out, err := dzce.Encode(kind, value)
if err != nil {
    return err
}

_ = out
```

Path-based load/save:

```go
kind, value, err := dzce.LoadFile("path/to/db/types.xml")
if err != nil {
    return err
}

_ = kind

if err := dzce.SaveFile("path/to/db/types.xml", value); err != nil {
    return err
}
```

## CE project Pipeline

Forward build:

```go
pipeline := dzce.CEPipeline{
    ConfigPath:       "ceproject/MyMap/mymap.xml",
    LayersDir:        "ceproject/MyMap/layers",
    TerritoriesDir:   "ceproject/MyMap/territoryTypes",
    AreaFlagsMapPath: "ceproject/MyMap/map/mymap.map",
}

if _, err := pipeline.Forward(); err != nil {
    return err
}
```

Reverse restore:

```go
pipeline := dzce.CEPipeline{
    ConfigPath:       "ceproject/MyMap/mymap.xml",
    LayersDir:        "ceproject/MyMap/layers",
    TerritoriesDir:   "ceproject/MyMap/territoryTypes",
    AreaFlagsMapPath: "ceproject/MyMap/map/mymap.map",
}

_, err := pipeline.Restore(dzce.CERestoreOptions{
    TemplateConfigPath: "backup/mymap.xml",
    LayerImageOptions: dzce.MaskImageEncodeOptions{
        Format:              dzce.MaskImageFormatTGA,
        ReuseBlankLayerFile: true,
        BlankLayerFileName:  "_blank",
    },
})
if err != nil {
    return err
}
```

Notes:

* map-only restore cannot reconstruct independent `usage=0/value=0` layers
* if those layers matter, keep source `layers/` as part of artifacts
* synthetic restore (without template) uses deterministic pseudo-random colors

## EconomyCore Merge

Merge `cfgeconomycore.xml` include tree:

```go
merged, err := dzce.LoadMergedEconomyCore("path/to/cfgeconomycore.xml")
if err != nil {
    return err
}

_ = merged
```

Relaxed include mode:

```go
merged, err := dzce.LoadMergedEconomyCoreWithOptions(
    "path/to/cfgeconomycore.xml",
    dzce.MergeOptions{RelaxedIncludeTypes: true},
)
if err != nil {
    return err
}

_ = merged
```
