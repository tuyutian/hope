export type User = {
  shop :string
  userGuide: UserGuide
}
export type UserGuide = {
  enabled:boolean
  setting_protension:boolean
  setup_widget:boolean
  how_work:boolean
  choose:boolean
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