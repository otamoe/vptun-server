import {RouteRecordRaw, createRouter as vueCreateRouter, createWebHistory, Router } from 'vue-router'

export function createRoutes(): RouteRecordRaw[] {
    const routes = [
        {
            path: '/',
            redirect: '/client',
        },
        {
            path: '/client',
            name: 'client/main',
            component: () => import('@/views/clients/Main.vue'),
            children: [
                {
                    path: '',
                    name: 'client/home',
                    component: () => import('@/views/clients/Home.vue'),
                },
                {
                    path: 'create',
                    name: 'client/create',
                    component: () => import('@/views/clients/Save.vue'),
                },
                {
                    path: ':client',
                    name: 'client/read',
                    component: () => import('@/views/clients/Read.vue'),
                },
                {
                    path: ':client/update',
                    name: 'client/update',
                    component: () => import('@/views/clients/Save.vue'),
                },
                {
                    path: ':client/shell',
                    name: 'client/shell/main',
                    component: () => import('@/views/clients/shell/Main.vue'),
                    children: [
                        {
                            path: '',
                            name: 'client/shell/home',
                            component: () => import('@/views/clients/shell/Home.vue'),
                        },
                        {
                            path: 'create',
                            name: 'client/shell/create',
                            component: () => import('@/views/clients/shell/Save.vue'),
                        },
                        {
                            path: ':shell',
                            name: 'client/shell/read',
                            component: () => import('@/views/clients/shell/Read.vue'),
                        },
                        {
                            path: ':shell/update',
                            name: 'client/shell/update',
                            component: () => import('@/views/clients/shell/Save.vue'),
                        },
                    ],
                },
            ]
        },
        {
            path: '/route',
            name: 'route/main',
            component: () => import('@/views/routes/Main.vue'),
            children: [
                {
                    path: '',
                    name: 'route/home',
                    component: () => import('@/views/routes/Home.vue'),
                },
                {
                    path: 'create',
                    name: 'route/create',
                    component: () => import('@/views/routes/Save.vue'),
                },
                {
                    path: ':route',
                    name: 'route/read',
                    component: () => import('@/views/routes/Read.vue'),
                },
                {
                    path: ':route/update',
                    name: 'route/update',
                    component: () => import('@/views/routes/Save.vue'),
                },
            ],
        },
        {
            path: '/:pathMatch(.*)*',
            name: 'not_found',
            component: () => import('@/views/pages/NotFound.vue'),
        },
    ]
    return routes
}

export function createRouter(): Router {
    const routes = createRoutes()
    return vueCreateRouter({routes, history: createWebHistory()})
}
