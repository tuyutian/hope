import React, { useState, useTransition } from "react";
import { Spinner } from "@shopify/polaris";
import { userService } from "@/api";
import { useAuth } from "@/stores/context.ts";
import { getUserState } from "@/stores/userStore.ts";
import { UserGuide } from "@/types/user.ts";

type Props = {
  name: keyof UserGuide;
  check: boolean;
};

export default function CheckCycle({ name, check }: Props) {
  const [done, setDone] = useState(check);
  const { user, setUser } = useAuth();
  const { toggleUserGuideStep } = getUserState();
  const [clickLoading, startTransition] = useTransition();
  const changeStatus = function (status: boolean) {
    startTransition(async () => {
      const res = await userService.updateGuide({ name, open: status });
      if (res.code === 0) {
        startTransition(() => {
          setDone(status);
          toggleUserGuideStep(name, status);
          setUser({ ...user });
        });
      }
    });
  };
  return clickLoading ? (
    <div className="flex items-center justify-center w-6 h-6 cursor-default">
      <Spinner accessibilityLabel="Small spinner example" size="small" />
    </div>
  ) : (
    <>
      {done ? (
        <div
          className="transition-transform w-5 h-5 m-[2px] transform cursor-pointer"
          onClick={() => {
            changeStatus(!done);
          }}
        >
          <svg width={20} height={20} viewBox="3 3 14 14" focusable="false" aria-hidden="true">
            <path d="M13.28 9.03a.75.75 0 0 0-1.06-1.06l-2.97 2.97-1.22-1.22a.75.75 0 0 0-1.06 1.06l1.75 1.75a.75.75 0 0 0 1.06 0l3.5-3.5Z" />
            <path
              fillRule="evenodd"
              d="M17 10a7 7 0 1 1-14 0 7 7 0 0 1 14 0Zm-1.5 0a5.5 5.5 0 1 1-11 0 5.5 5.5 0 0 1 11 0Z"
            />
          </svg>
        </div>
      ) : (
        <div
          className="w-5 min-w-5 h-5 m-[2px] rounded-[30px] border-2 border-dashed hover:border-solid cursor-pointer active:bg-slate-50"
          style={{ borderColor: "#8A8A8A" }}
          onClick={e => {
            e.stopPropagation();
            changeStatus(!done);
          }}
        />
      )}
    </>
  );
}
