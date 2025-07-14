import React from "react";
import { Tooltip } from "@shopify/polaris";

interface TooltipsProps {
  width: number;
  title: string;
  text: React.ReactNode;
}

const HomeTooltips: React.FC<TooltipsProps> = ({ title, text }) => {
  return (
    <div>
      <Tooltip content={text}>
        <p className="leading-5 underline decoration-dashed decoration-gray-500">{title}</p>
      </Tooltip>
    </div>
  );
};

export default HomeTooltips;
