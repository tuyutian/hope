export type User = {
  shop :string
  userGuide:Record<string, boolean>
}
export const DefaultUser:User = {
  shop: "",
  userGuide: {
    "enabled": false,
    "setting_protension":false,
    "setup_widget":false,
    "how_work":false,
    "choose":false,
  },
}