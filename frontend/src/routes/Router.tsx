import {createBrowserRouter, RouteObject, RouterProvider} from "react-router";
import React, {lazy, Suspense} from "react";

// 布局组件
const Layout = lazy(() => import("@/layouts/MainLayout"));

// 页面组件 - 使用懒加载优化性能
const Home = lazy(() => import("@/pages/Home"));
const Order = lazy(() => import("@/pages/Order"));
const Cart = lazy(() => import("@/pages/Cart"));
const NotFound = lazy(() => import("@/pages/NotFound"));

// 加载指示器组件
const LoadingFallback = () => <div className="loading-spinner">加载中...</div>;

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
        path: "order",
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