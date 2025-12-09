# MicroSQL AGo - Frontend

Frontend de la aplicación MicroSQL AGo construido con Next.js 16, React 19, TypeScript y Tailwind CSS.

## 🚀 Tecnologías

- **Next.js 16** - Framework React con App Router
- **React 19** - Biblioteca de UI
- **TypeScript** - Tipado estático
- **Tailwind CSS 4** - Framework de estilos utility-first
- **ESLint** - Linter para código JavaScript/TypeScript

## 📁 Estructura del Proyecto

```
frontend-nextjs/
├── src/
│   └── app/              # App Router de Next.js
│       ├── layout.tsx    # Layout raíz
│       ├── page.tsx      # Página principal
│       └── globals.css   # Estilos globales
├── public/               # Archivos estáticos
├── next.config.ts        # Configuración de Next.js
├── tsconfig.json         # Configuración de TypeScript
└── package.json          # Dependencias del proyecto
```

## 🛠️ Instalación

Asegúrate de tener Node.js instalado (versión 18 o superior).

```bash
# Instalar dependencias
npm install

# O con yarn
yarn install

# O con pnpm
pnpm install
```

## 🏃 Desarrollo

Inicia el servidor de desarrollo:

```bash
npm run dev
# o
yarn dev
# o
pnpm dev
```

Abre [http://localhost:3000](http://localhost:3000) en tu navegador para ver la aplicación.

La página se actualiza automáticamente cuando editas los archivos.

## 📦 Scripts Disponibles

- `npm run dev` - Inicia el servidor de desarrollo
- `npm run build` - Construye la aplicación para producción
- `npm run start` - Inicia el servidor de producción (después de `build`)
- `npm run lint` - Ejecuta ESLint para verificar el código

## 🔗 Integración con Backend

Este frontend se conecta con el backend Go ubicado en `../backend-go/`.

**URL del Backend (desarrollo):** `http://localhost:8080`

### Endpoints principales:
- `/api/auth/login` - Autenticación de usuarios
- `/api/auth/logout` - Cerrar sesión
- `/api/users/register` - Registro de usuarios
- `/api/db/*` - Gestión de conexiones a bases de datos
- `/api/admin/*` - Endpoints de administración

## 📝 Próximos Pasos

1. Configurar variables de entorno para la URL del backend
2. Crear servicios de API para comunicarse con el backend
3. Implementar autenticación y manejo de tokens JWT
4. Crear componentes reutilizables
5. Implementar las páginas principales:
   - Login/Registro
   - Dashboard
   - Gestión de conexiones
   - Ejecución de auditorías
   - Panel de administración

## 🎨 Estilos

El proyecto usa Tailwind CSS 4 con configuración moderna. Los estilos globales están en `src/app/globals.css`.

## 📚 Recursos

- [Documentación de Next.js](https://nextjs.org/docs)
- [Documentación de React](https://react.dev)
- [Documentación de Tailwind CSS](https://tailwindcss.com/docs)
- [Documentación de TypeScript](https://www.typescriptlang.org/docs)

## 🚢 Despliegue

La forma más fácil de desplegar tu aplicación Next.js es usando [Vercel Platform](https://vercel.com/new).

Consulta la [documentación de despliegue de Next.js](https://nextjs.org/docs/app/building-your-application/deploying) para más detalles.
