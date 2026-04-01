// SPDX-License-Identifier: MIT
// Copyright (c) 2026 WoozyMasta
// Source: github.com/woozymasta/dzce

package dzce

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/woozymasta/bimime"
)

// defaultRegistry stores built-in CE formats.
var defaultRegistry = NewRegistry(
	newXMLFormat[TypesFile](KindTypes),
	newXMLFormat[EventsFile](KindEvents),
	newXMLFormat[EconomyFile](KindEconomy),
	newXMLFormat[GlobalsFile](KindGlobals),
	newXMLFormat[MessagesFile](KindMessages),
	newXMLFormat[SpawnableTypesFile](KindSpawnableTypes),
	newXMLFormat[RandomPresetsFile](KindRandomPresets),
	newXMLFormat[EconomyCoreFile](KindEconomyCore),
	newXMLFormat[EnvironmentFile](KindEnvironment),
	newXMLFormat[EventSpawnsFile](KindEventSpawns),
	newXMLFormat[EventGroupsFile](KindEventGroups),
	newXMLFormat[PlayerSpawnPointsFile](KindPlayerSpawnPoints),
	newXMLFormat[WeatherFile](KindWeather),
	newXMLFormat[LimitsDefinitionFile](KindLimitsDefinition),
	newXMLFormat[LimitsDefinitionUserFile](KindLimitsDefinitionUser),
	newXMLFormat[IgnoreListFile](KindIgnoreList),
	newXMLFormat[TerritoryFile](KindTerritories),
	newXMLFormat[MapGroupProtoFile](KindMapGroupProto),
	newXMLFormat[MapClusterProtoFile](KindMapClusterProto),
	newXMLFormat[MapGroupPosFile](KindMapGroupPos),
	newXMLFormat[MapGroupDirtFile](KindMapGroupDirt),
	newXMLFormat[MapGroupClusterFile](KindMapGroupCluster),
	newXMLFormat[CEProjectConfigFile](KindCEProjectConfig),
	newJSONFormat[UndergroundTriggersFile](KindUndergroundTriggers),
	newJSONFormat[EffectAreaFile](KindEffectArea),
	newJSONFormat[GameplayFile](KindGameplay),
	newJSONFormat[GameplayGearPresetsFile](KindGameplayGearPresets),
	newJSONFormat[ObjectSpawnerFile](KindObjectSpawner),
	newAreaFlagsMapFormat(),
)

// Format provides decode/encode operations for one CE kind.
type Format interface {
	Kind() Kind
	Decode(data []byte) (any, error)
	Encode(value any) ([]byte, error)
}

// Registry provides CE kind-to-format resolution.
type Registry interface {
	Get(kind Kind) (Format, bool)
	Decode(kind Kind, data []byte) (any, error)
	Encode(kind Kind, value any) ([]byte, error)
}

// MapRegistry is a map-based Registry implementation.
type MapRegistry struct {
	formats map[Kind]Format
}

// NewRegistry creates a map registry from provided formats.
func NewRegistry(formats ...Format) *MapRegistry {
	items := make(map[Kind]Format, len(formats))

	for _, format := range formats {
		items[format.Kind()] = format
	}

	return &MapRegistry{formats: items}
}

// Get returns a registered format by CE kind.
func (r *MapRegistry) Get(kind Kind) (Format, bool) {
	format, ok := r.formats[kind]
	return format, ok
}

// Decode decodes CE payload for provided kind.
func (r *MapRegistry) Decode(kind Kind, data []byte) (any, error) {
	format, ok := r.Get(kind)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedKind, kind)
	}

	value, err := format.Decode(data)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Encode encodes CE payload for provided kind.
func (r *MapRegistry) Encode(kind Kind, value any) ([]byte, error) {
	format, ok := r.Get(kind)
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedKind, kind)
	}

	data, err := format.Encode(value)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Decode decodes CE payload with default registry.
func Decode(kind Kind, data []byte) (any, error) {
	return defaultRegistry.Decode(kind, data)
}

// DecodeReader decodes CE payload with default registry from reader.
func DecodeReader(kind Kind, reader io.Reader) (any, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read ce content: %w", err)
	}

	return Decode(kind, data)
}

// Encode encodes CE payload with default registry.
func Encode(kind Kind, value any) ([]byte, error) {
	return defaultRegistry.Encode(kind, value)
}

// EncodeWriter encodes CE payload with default registry into writer.
func EncodeWriter(kind Kind, writer io.Writer, value any) error {
	data, err := Encode(kind, value)
	if err != nil {
		return err
	}

	if _, err = writer.Write(data); err != nil {
		return fmt.Errorf("write ce content: %w", err)
	}

	return nil
}

// DefaultRegistry returns built-in CE format registry.
func DefaultRegistry() Registry {
	return defaultRegistry
}

// LoadFile loads CE file by detected kind.
func LoadFile(path string) (Kind, any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return KindUnknown, nil, fmt.Errorf("read ce file %q: %w", path, err)
	}

	kind := DetectKind(path)
	if kind == KindUnknown {
		kind = detectKindByContent(path, data)
	}

	if kind == KindUnknown {
		return KindUnknown, nil, fmt.Errorf("%w: %s", ErrUnknownFileKind, path)
	}

	value, err := Decode(kind, data)
	if err != nil {
		return KindUnknown, nil, fmt.Errorf("decode ce file %q: %w", path, err)
	}

	return kind, value, nil
}

// LoadFileAs loads CE file using explicit kind.
//
// This is useful for JSON files with arbitrary names, such as
// `spawnGearPresetFiles` and `objectSpawnersArr` entries.
func LoadFileAs(kind Kind, path string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read ce file %q: %w", path, err)
	}

	value, err := Decode(kind, data)
	if err != nil {
		return nil, fmt.Errorf("decode ce file %q: %w", path, err)
	}

	return value, nil
}

// SaveFile stores CE file by detected kind.
func SaveFile(path string, value any) error {
	kind := DetectKind(path)
	if kind == KindUnknown {
		kind = detectKindByValue(value)
	}

	if kind == KindUnknown {
		return fmt.Errorf("%w: %s", ErrUnknownFileKind, path)
	}

	data, err := Encode(kind, value)
	if err != nil {
		return fmt.Errorf("encode ce file %q: %w", path, err)
	}

	if err := writeFile600(path, data); err != nil {
		return fmt.Errorf("write ce file %q: %w", path, err)
	}

	return nil
}

// detectKindByContent detects known dynamic-name CE formats by payload.
func detectKindByContent(path string, data []byte) Kind {
	result := bimime.Analyze(bimime.AnalyzeOptions{
		Path:        path,
		Prefix:      data,
		DefaultPlan: bimime.PlanNormal(),
	})

	return kindFromTypeID(result.Probe.Resolved.ID)
}

// detectKindByValue detects known dynamic-name CE formats by value type.
func detectKindByValue(value any) Kind {
	switch value.(type) {
	case *CEProjectConfigFile:
		return KindCEProjectConfig
	case *AreaFlagsMapFile:
		return KindAreaFlagsMap
	default:
		return KindUnknown
	}
}

// SaveFileAs stores CE file using explicit kind.
//
// This is useful for JSON files with arbitrary names, such as
// `spawnGearPresetFiles` and `objectSpawnersArr` entries.
func SaveFileAs(kind Kind, path string, value any) error {
	data, err := Encode(kind, value)
	if err != nil {
		return fmt.Errorf("encode ce file %q: %w", path, err)
	}

	if err := writeFile600(path, data); err != nil {
		return fmt.Errorf("write ce file %q: %w", path, err)
	}

	return nil
}

// DecodeTypes decodes `db/types.xml`.
func DecodeTypes(data []byte) (*TypesFile, error) {
	return decodeXML[TypesFile](data)
}

// EncodeTypes encodes `db/types.xml`.
func EncodeTypes(value *TypesFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeEvents decodes `db/events.xml`.
func DecodeEvents(data []byte) (*EventsFile, error) {
	return decodeXML[EventsFile](data)
}

// EncodeEvents encodes `db/events.xml`.
func EncodeEvents(value *EventsFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeEconomy decodes `db/economy.xml`.
func DecodeEconomy(data []byte) (*EconomyFile, error) {
	return decodeXML[EconomyFile](data)
}

// EncodeEconomy encodes `db/economy.xml`.
func EncodeEconomy(value *EconomyFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeGlobals decodes `db/globals.xml`.
func DecodeGlobals(data []byte) (*GlobalsFile, error) {
	return decodeXML[GlobalsFile](data)
}

// EncodeGlobals encodes `db/globals.xml`.
func EncodeGlobals(value *GlobalsFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeMessages decodes `db/messages.xml`.
func DecodeMessages(data []byte) (*MessagesFile, error) {
	return decodeXML[MessagesFile](data)
}

// EncodeMessages encodes `db/messages.xml`.
func EncodeMessages(value *MessagesFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeSpawnableTypes decodes `cfgspawnabletypes.xml`.
func DecodeSpawnableTypes(data []byte) (*SpawnableTypesFile, error) {
	return decodeXML[SpawnableTypesFile](data)
}

// EncodeSpawnableTypes encodes `cfgspawnabletypes.xml`.
func EncodeSpawnableTypes(value *SpawnableTypesFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeRandomPresets decodes `cfgrandompresets.xml`.
func DecodeRandomPresets(data []byte) (*RandomPresetsFile, error) {
	return decodeXML[RandomPresetsFile](data)
}

// EncodeRandomPresets encodes `cfgrandompresets.xml`.
func EncodeRandomPresets(value *RandomPresetsFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeEconomyCore decodes `cfgeconomycore.xml`.
func DecodeEconomyCore(data []byte) (*EconomyCoreFile, error) {
	return decodeXML[EconomyCoreFile](data)
}

// EncodeEconomyCore encodes `cfgeconomycore.xml`.
func EncodeEconomyCore(value *EconomyCoreFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeEnvironment decodes `cfgenvironment.xml`.
func DecodeEnvironment(data []byte) (*EnvironmentFile, error) {
	return decodeXML[EnvironmentFile](data)
}

// EncodeEnvironment encodes `cfgenvironment.xml`.
func EncodeEnvironment(value *EnvironmentFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeEventSpawns decodes `cfgeventspawns.xml`.
func DecodeEventSpawns(data []byte) (*EventSpawnsFile, error) {
	return decodeXML[EventSpawnsFile](data)
}

// EncodeEventSpawns encodes `cfgeventspawns.xml`.
func EncodeEventSpawns(value *EventSpawnsFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeEventGroups decodes `cfgeventgroups.xml`.
func DecodeEventGroups(data []byte) (*EventGroupsFile, error) {
	return decodeXML[EventGroupsFile](data)
}

// EncodeEventGroups encodes `cfgeventgroups.xml`.
func EncodeEventGroups(value *EventGroupsFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodePlayerSpawnPoints decodes `cfgplayerspawnpoints.xml`.
func DecodePlayerSpawnPoints(data []byte) (*PlayerSpawnPointsFile, error) {
	return decodeXML[PlayerSpawnPointsFile](data)
}

// EncodePlayerSpawnPoints encodes `cfgplayerspawnpoints.xml`.
func EncodePlayerSpawnPoints(value *PlayerSpawnPointsFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeWeather decodes `cfgweather.xml`.
func DecodeWeather(data []byte) (*WeatherFile, error) {
	return decodeXML[WeatherFile](data)
}

// EncodeWeather encodes `cfgweather.xml`.
func EncodeWeather(value *WeatherFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeLimitsDefinition decodes `cfglimitsdefinition.xml`.
func DecodeLimitsDefinition(data []byte) (*LimitsDefinitionFile, error) {
	return decodeXML[LimitsDefinitionFile](data)
}

// EncodeLimitsDefinition encodes `cfglimitsdefinition.xml`.
func EncodeLimitsDefinition(value *LimitsDefinitionFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeLimitsDefinitionUser decodes `cfglimitsdefinitionuser.xml`.
func DecodeLimitsDefinitionUser(data []byte) (*LimitsDefinitionUserFile, error) {
	return decodeXML[LimitsDefinitionUserFile](data)
}

// EncodeLimitsDefinitionUser encodes `cfglimitsdefinitionuser.xml`.
func EncodeLimitsDefinitionUser(value *LimitsDefinitionUserFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeIgnoreList decodes `cfgignorelist.xml`.
func DecodeIgnoreList(data []byte) (*IgnoreListFile, error) {
	return decodeXML[IgnoreListFile](data)
}

// EncodeIgnoreList encodes `cfgignorelist.xml`.
func EncodeIgnoreList(value *IgnoreListFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeTerritory decodes `env/*_territories.xml`.
func DecodeTerritory(data []byte) (*TerritoryFile, error) {
	return decodeXML[TerritoryFile](data)
}

// EncodeTerritory encodes `env/*_territories.xml`.
func EncodeTerritory(value *TerritoryFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeMapGroupProto decodes `mapgroupproto.xml`.
func DecodeMapGroupProto(data []byte) (*MapGroupProtoFile, error) {
	return decodeXML[MapGroupProtoFile](data)
}

// EncodeMapGroupProto encodes `mapgroupproto.xml`.
func EncodeMapGroupProto(value *MapGroupProtoFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeMapClusterProto decodes `mapclusterproto.xml`.
func DecodeMapClusterProto(data []byte) (*MapClusterProtoFile, error) {
	return decodeXML[MapClusterProtoFile](data)
}

// EncodeMapClusterProto encodes `mapclusterproto.xml`.
func EncodeMapClusterProto(value *MapClusterProtoFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeMapGroupPos decodes `mapgrouppos.xml`.
func DecodeMapGroupPos(data []byte) (*MapGroupPosFile, error) {
	return decodeXML[MapGroupPosFile](data)
}

// EncodeMapGroupPos encodes `mapgrouppos.xml`.
func EncodeMapGroupPos(value *MapGroupPosFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeMapGroupDirt decodes `mapgroupdirt.xml`.
func DecodeMapGroupDirt(data []byte) (*MapGroupDirtFile, error) {
	return decodeXML[MapGroupDirtFile](data)
}

// EncodeMapGroupDirt encodes `mapgroupdirt.xml`.
func EncodeMapGroupDirt(value *MapGroupDirtFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeMapGroupCluster decodes `mapgroupcluster*.xml`.
func DecodeMapGroupCluster(data []byte) (*MapGroupClusterFile, error) {
	return decodeXML[MapGroupClusterFile](data)
}

// EncodeMapGroupCluster encodes `mapgroupcluster*.xml`.
func EncodeMapGroupCluster(value *MapGroupClusterFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeCEProjectConfig decodes CEProject `mapname.xml` (`<zg-config>`).
func DecodeCEProjectConfig(data []byte) (*CEProjectConfigFile, error) {
	return decodeXML[CEProjectConfigFile](data)
}

// EncodeCEProjectConfig encodes CEProject `mapname.xml` (`<zg-config>`).
func EncodeCEProjectConfig(value *CEProjectConfigFile) ([]byte, error) {
	return encodeXML(value)
}

// DecodeUndergroundTriggers decodes `cfgundergroundtriggers.json`.
func DecodeUndergroundTriggers(data []byte) (*UndergroundTriggersFile, error) {
	return decodeJSON[UndergroundTriggersFile](data)
}

// EncodeUndergroundTriggers encodes `cfgundergroundtriggers.json`.
func EncodeUndergroundTriggers(value *UndergroundTriggersFile) ([]byte, error) {
	return encodeJSON(value)
}

// DecodeEffectArea decodes `cfgeffectarea.json`.
func DecodeEffectArea(data []byte) (*EffectAreaFile, error) {
	return decodeJSON[EffectAreaFile](data)
}

// EncodeEffectArea encodes `cfgeffectarea.json`.
func EncodeEffectArea(value *EffectAreaFile) ([]byte, error) {
	return encodeJSON(value)
}

// DecodeGameplay decodes `cfggameplay.json`.
func DecodeGameplay(data []byte) (*GameplayFile, error) {
	return decodeJSON[GameplayFile](data)
}

// EncodeGameplay encodes `cfggameplay.json`.
func EncodeGameplay(value *GameplayFile) ([]byte, error) {
	return encodeJSON(value)
}

// DecodeGameplayGearPresets decodes one gear preset JSON payload.
//
// The file name is arbitrary and usually comes from
// `cfggameplay.json -> PlayerData.spawnGearPresetFiles`.
func DecodeGameplayGearPresets(data []byte) (*GameplayGearPresetsFile, error) {
	return decodeJSON[GameplayGearPresetsFile](data)
}

// EncodeGameplayGearPresets encodes one gear preset JSON payload.
func EncodeGameplayGearPresets(value *GameplayGearPresetsFile) ([]byte, error) {
	return encodeJSON(value)
}

// DecodeObjectSpawner decodes one object spawner JSON payload.
//
// The file name is arbitrary and usually comes from
// `cfggameplay.json -> WorldsData.objectSpawnersArr`.
func DecodeObjectSpawner(data []byte) (*ObjectSpawnerFile, error) {
	return decodeJSON[ObjectSpawnerFile](data)
}

// EncodeObjectSpawner encodes one object spawner JSON payload.
func EncodeObjectSpawner(value *ObjectSpawnerFile) ([]byte, error) {
	return encodeJSON(value)
}

// areaFlagsMapFormat decodes and encodes `areaflags.map` binary payload.
type areaFlagsMapFormat struct{}

// newAreaFlagsMapFormat creates `areaflags.map` format implementation.
func newAreaFlagsMapFormat() areaFlagsMapFormat {
	return areaFlagsMapFormat{}
}

// Kind returns CE kind for `areaflags.map`.
func (areaFlagsMapFormat) Kind() Kind {
	return KindAreaFlagsMap
}

// Decode parses `areaflags.map` payload.
func (areaFlagsMapFormat) Decode(data []byte) (any, error) {
	value, err := DecodeAreaFlagsMap(data)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", KindAreaFlagsMap, err)
	}

	return value, nil
}

// Encode serializes `areaflags.map` payload.
func (areaFlagsMapFormat) Encode(value any) ([]byte, error) {
	target, ok := castValue[AreaFlagsMapFile](value)
	if !ok || target == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, KindAreaFlagsMap)
	}

	data, err := EncodeAreaFlagsMap(target)
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", KindAreaFlagsMap, err)
	}

	return data, nil
}

// jsonFormat decodes and encodes JSON payload for one root type.
type jsonFormat[T any] struct {
	kind Kind
}

// newJSONFormat creates JSON format implementation for a CE kind.
func newJSONFormat[T any](kind Kind) jsonFormat[T] {
	return jsonFormat[T]{kind: kind}
}

// Kind returns CE kind for the format.
func (f jsonFormat[T]) Kind() Kind {
	return f.kind
}

// Decode parses JSON payload into typed value.
func (f jsonFormat[T]) Decode(data []byte) (any, error) {
	value, err := decodeJSON[T](data)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", f.kind, err)
	}

	return value, nil
}

// Encode serializes typed value into JSON payload.
func (f jsonFormat[T]) Encode(value any) ([]byte, error) {
	target, ok := castValue[T](value)
	if !ok || target == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, f.kind)
	}

	data, err := encodeJSON(target)
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", f.kind, err)
	}

	return data, nil
}

// xmlFormat decodes and encodes XML payload for one root type.
type xmlFormat[T any] struct {
	kind Kind
}

// newXMLFormat creates XML format implementation for a CE kind.
func newXMLFormat[T any](kind Kind) xmlFormat[T] {
	return xmlFormat[T]{kind: kind}
}

// Kind returns CE kind for the format.
func (f xmlFormat[T]) Kind() Kind {
	return f.kind
}

// Decode parses XML payload into typed value.
func (f xmlFormat[T]) Decode(data []byte) (any, error) {
	value, err := decodeXML[T](data)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %w", f.kind, err)
	}

	return value, nil
}

// Encode serializes typed value into XML payload.
func (f xmlFormat[T]) Encode(value any) ([]byte, error) {
	target, ok := castValue[T](value)
	if !ok || target == nil {
		return nil, fmt.Errorf("%w for %s", ErrUnsupportedValue, f.kind)
	}

	data, err := encodeXML(target)
	if err != nil {
		return nil, fmt.Errorf("encode %s: %w", f.kind, err)
	}

	return data, nil
}

// decodeXML decodes XML payload into type T.
func decodeXML[T any](data []byte) (*T, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	var value T

	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return &value, nil
}

// decodeJSON decodes JSON payload into type T.
func decodeJSON[T any](data []byte) (*T, error) {
	var value T

	if err := json.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	return &value, nil
}

// encodeJSON encodes type T into indented JSON payload.
func encodeJSON[T any](value *T) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("%w: nil", ErrUnsupportedValue)
	}

	data, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("encode json: %w", err)
	}

	if len(data) == 0 || data[len(data)-1] != '\n' {
		data = append(data, '\n')
	}

	return data, nil
}

// encodeXML encodes type T into XML payload with header and indentation.
func encodeXML[T any](value *T) ([]byte, error) {
	if value == nil {
		return nil, fmt.Errorf("%w: nil", ErrUnsupportedValue)
	}

	buffer := bytes.NewBufferString(xml.Header)
	encoder := xml.NewEncoder(buffer)
	encoder.Indent("", "    ")

	if err := encoder.Encode(value); err != nil {
		return nil, fmt.Errorf("encode xml: %w", err)
	}

	if err := encoder.Flush(); err != nil {
		return nil, fmt.Errorf("flush xml: %w", err)
	}

	if buffer.Len() == 0 || buffer.Bytes()[buffer.Len()-1] != '\n' {
		buffer.WriteByte('\n')
	}

	return buffer.Bytes(), nil
}

// castValue casts incoming value to `*T` or creates pointer from `T`.
func castValue[T any](value any) (*T, bool) {
	if pointer, ok := value.(*T); ok {
		return pointer, true
	}

	if direct, ok := value.(T); ok {
		return &direct, true
	}

	return nil, false
}
