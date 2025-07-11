import React, {useState} from "react";
import {Spinner} from "@shopify/polaris";

export default function CheckCycle() {
  const [clickLoading, setClickLoading] = useState(false);
  const [done, setDone] = useState(false);
  const [atChange, setAtChange] = useState(false);
  const changeStatus = function (status: boolean) {
    setDone(status);
  };
  return clickLoading ? <div className="flex items-center justify-center w-6 h-6 cursor-default">
    <Spinner accessibilityLabel="Small spinner example" size="small" />
  </div> : <>
    {done ? (
      <div
        className={`transition-transform w-5 h-5 m-[2px] transform ${atChange ? "cursor-not-allowed" : "cursor-pointer"}`}
        onClick={e => {
          e.stopPropagation();
          if (atChange) return;
          changeStatus(!done);
        }}
      >
        <svg width={20} height={20} viewBox="3 3 14 14" focusable="false" aria-hidden="true">
          <path
            d="M13.28 9.03a.75.75 0 0 0-1.06-1.06l-2.97 2.97-1.22-1.22a.75.75 0 0 0-1.06 1.06l1.75 1.75a.75.75 0 0 0 1.06 0l3.5-3.5Z" />
          <path fillRule="evenodd"
                d="M17 10a7 7 0 1 1-14 0 7 7 0 0 1 14 0Zm-1.5 0a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0Z" />
        </svg>
      </div>
    ) : (
      <div
        className={`w-5 min-w-5 h-5 m-[2px] rounded-[30px] border-2 border-dashed hover:border-solid ${
          atChange ? "cursor-not-allowed" : "cursor-pointer"
        } active:bg-slate-50`}
        style={{borderColor: "#8A8A8A"}}
        onClick={e => {
          e.stopPropagation();
          if (atChange) return;
          changeStatus(!done);
        }}
      />
    )}
  </>;
}
