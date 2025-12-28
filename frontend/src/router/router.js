import { createWebHistory, createRouter } from "vue-router";

// определяем маршруты
const routes = [
    {
        path: "/",
        name: "home",
        component: () => import('../views/HomeView.vue'),
        props: true, // указываем, что компонент Category.vue может принимать параметры в адресной строке, например, в path указан id
        meta: {
            title: "Добро пожаловать!"
        }
    },
    {
    path: "/main",
        name: "main",
        component: () => import('../views/Categories.vue'),
        props: true, // указываем, что компонент Category.vue может принимать параметры в адресной строке, например, в path указан id
        meta: {
            title: "Квизы!"
        }
    },
    {
        path: "/login", // указание маршрута, по которому будем переходить к списку категорий
        name: "login", // имя маршрута
        component: () => import('../views/Login.vue'), // компонент, на основании которого будет отрисовываться страница
        meta: {
            title: "Вход"
        }
    },
    {
        path: "/register",
        name: "register",
        component: () => import('../views/Register.vue'),
        props: true, // указываем, что компонент Category.vue может принимать параметры в адресной строке, например, в path указан id
        meta: {
            title: "Регистрация"
        }
    },
    {
        path: '/category/:id/quizzes',
        name: 'CategoryQuizzes',
        component: () => import('../views/CategoryQuizzes.vue'),
        props: true,
        meta: {
            title: "Квизы"
        }
    },
    {
        path: '/quiz/:id',
        name: 'quiz',
        component: () => import('../views/Quiz.vue'),
        props: true,
        meta: {
            title: "Квиз"
        }
    },
    {
        path: '/quiz/:id/play',
        name: 'quizgo',
        component: () => import('../views/Question.vue'),
        props: true,
        meta: {
            title: "Прохождение квиза"
        }
    },
    {
        path: "/create-quiz",
        name: "create",
        component: () => import('../views/QuizAdd.vue'),
        props: true, // указываем, что компонент Category.vue может принимать параметры в адресной строке, например, в path указан id
        meta: {
            title: "Создание квиза"
        }
    }
];

const router = createRouter({
    history: createWebHistory(), // указываем, что будет создаваться история посещений веб-страниц
    routes, // подключаем маршрутизацию
});

// указание заголовка компонентам (тега title), заголовки определены в meta
router.beforeEach((to, from, next) => {
    // для тех маршрутов, для которых не определены компоненты, подключается только App.vue
    // поэтому устанавливаем заголовком по умолчанию название "Главная страница"
    document.title = to.meta.title || 'Главная страница';
    next();
});

export default router;