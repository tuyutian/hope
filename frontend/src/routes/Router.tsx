import {createBrowserRouter, RouteObject, RouterProvider} from "react-router";
import React, {lazy, Suspense} from "react";
import classNames from 'classnames'

// 布局组件
const Layout = lazy(() => import("@/layouts/MainLayout"));

// 页面组件 - 使用懒加载优化性能
const Home = lazy(() => import("@/pages/Home"));
const Order = lazy(() => import("@/pages/Order"));
const Cart = lazy(() => import("@/pages/Cart"));
const NotFound = lazy(() => import("@/pages/NotFound"));
const Billing = lazy(() => import("@/pages/Billing"));
// 加载指示器组件
export const LoadingFallback = () => {
  return <div
        className="flex items-center justify-center absolute top-0 bg-white bg-opacity-90 w-full h-full"
        style={{
          zIndex: 98,
        }}
      >
        <s-spinner accessibilityLabel="Loading" size="large-100" />
      </div>
}

// 路由配置对象
const routes: RouteObject[] = [
  {
    path: "/",
    element: (
      <Suspense fallback={<LoadingFallback />}>
        <Layout />
      </Suspense>
    ),
    children: [
      {
        index: true,
        element: (
          <Suspense fallback={<LoadingFallback />}>
            <Home />
          </Suspense>
        ),
      },
      {
        path: "orders",
        element: (
          <Suspense fallback={<LoadingFallback />}>
            <Order />
          </Suspense>
        ),
      },
      {
        path: "cart",
        element: (
          <Suspense fallback={<LoadingFallback />}>
            <Cart />
          </Suspense>
        ),
      },

      {
        path: "billing",
        element: (
          <Suspense fallback={<LoadingFallback />}>
            <Billing />
          </Suspense>
        ),
      },
      {
        path: "*",
        element: (
          <Suspense fallback={<LoadingFallback />}>
            <NotFound />
          </Suspense>
        ),
      },
    ],
  },
];

// 创建路由器
const router = createBrowserRouter(routes);
const Router = () => {
  return <RouterProvider router={router}  />;
};

export default Router;