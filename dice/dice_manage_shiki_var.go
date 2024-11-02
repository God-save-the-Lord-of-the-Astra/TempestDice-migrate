package dice

import "fmt"

type ShikiGroupConfig struct {
	GroupID string
	Type    string
	Value   any
}

func (d *Dice) shikiSetGroupConfig(groupID string, key string, luaType string, value any) {
	val := ShikiGroupConfig{
		GroupID: groupID,
		Type:    luaType,
		Value:   value,
	}
	d.ConfigManager.SetConfig(fmt.Sprintf("shiki_group_%s", groupID), key, val)
}

func (d *Dice) shikiGetGroupConfig(groupID string, key string) (luaType string, value any) {
	val := d.ConfigManager.getConfig(fmt.Sprintf("shiki_group_%s", groupID), key)
	if val == nil {
		return "", nil
	}
	return val.Type, val.Value
}

func (d *Dice) shikiSetUserConfig(userID string, key string, luaType string, value any) {
	val := ShikiGroupConfig{
		GroupID: userID,
		Type:    luaType,
		Value:   value,
	}
	d.ConfigManager.SetConfig(fmt.Sprintf("shiki_user_%s", userID), key, val)
}

func (d *Dice) shikiGetUserConfig(userID string, key string) (luaType string, value any) {
	val := d.ConfigManager.getConfig(fmt.Sprintf("shiki_user_%s", userID), key)
	if val == nil {
		return "", nil
	}
	return val.Type, val.Value
}
