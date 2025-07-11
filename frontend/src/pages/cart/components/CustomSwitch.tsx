import React from "react";

interface SwitchProps {
  onChange: (checked: boolean) => void;
  checked: boolean;
  onColor?: string;
  offColor?: string;
  uncheckedIcon?: React.ReactNode;
  checkedIcon?: React.ReactNode | false;
}

const CustomSwitch: React.FC<SwitchProps> = ({
  onChange,
  checked,
  onColor = "",
  offColor = "",
  uncheckedIcon,
  checkedIcon
}) => (
  <div
    className={`inline-flex items-center justify-center min-w-8 w-8 h-5 rounded-md cursor-pointer transition-colors ${
      checked ? "bg-[#303030]" : "bg-[#E3E3E3]"
    }`}
    style={{backgroundColor: checked ? onColor : offColor}}
    onClick={() => onChange(!checked)}
  >
    <div
      className={`w-4 h-4 bg-white rounded-full shadow-md transform transition-transform ${
        checked ? "translate-x-1.5" : "-translate-x-1.5"
      }`}
    >
      {checked ? checkedIcon : uncheckedIcon}
    </div>
  </div>
);

export default CustomSwitch;
