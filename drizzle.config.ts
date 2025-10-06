// Drizzle configuration for database management
// This file is used by Drizzle Studio running in Docker container
// The drizzle-kit dependency is installed automatically in the container

export default {
  schema: "./internal/domain/*.go", // Path ke Go domain files
  out: "./migrations",
  dialect: "postgresql",
  dbCredentials: {
    host: "localhost",
    port: 5433,
    user: "4dm1n",
    password: "4dm1nP@ssw0rd",
    database: "dashboard_ac_dev",
  },
  // Konfigurasi untuk introspection dari database yang sudah ada
  introspect: {
    casing: "snake_case",
  },
};