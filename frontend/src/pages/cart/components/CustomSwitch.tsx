import React from "react";

interface SwitchProps {
  onChange: (checked: boolean) => void;
  checked: boolean;
  onColor?: string;
  offColor?: string;
  uncheckedIcon?: React.ReactNode;
  checkedIcon?: React.ReactNode;
  loading?: boolean;
}

const CustomSwitch: React.FC<SwitchProps> = ({
  onChange,
  checked,
  onColor = "",
  offColor = "",
  uncheckedIcon,
  checkedIcon,
  loading = false,
}) => (
  <div
    className={`inline-flex items-center justify-center min-w-8 w-8 h-5 rounded-md transition-colors ${
      loading ? "cursor-not-allowed opacity-60" : "cursor-pointer hover:opacity-90"
    } ${checked ? "bg-[#303030]" : "bg-[#E3E3E3]"}`}
    style={{ backgroundColor: checked ? onColor : offColor }}
    onClick={() => {
      if (!loading) {
        onChange(!checked);
      }
    }}
  >
    <div
      className={`w-4 h-4 bg-white rounded-full shadow-md transform transition-transform ${
        checked ? "translate-x-1.5" : "-translate-x-1.5"
      }`}
    >
      {loading ? (
        <div className="flex items-center justify-center w-full h-full">
          <div className="w-4 h-4 border border-gray-400 border-t-transparent rounded-full animate-spin" />
        </div>
      ) : checked ? (
        checkedIcon
      ) : (
        uncheckedIcon
      )}
    </div>
  </div>
);

export default CustomSwitch;
