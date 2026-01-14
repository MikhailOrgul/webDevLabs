<template>
  <main class="container mt-4">
    <h1>О компании &lt;Авто-Оргул&gt;</h1>

    <p>
      Четкие тазы для четких парней. У нас есть все тачки пачки катушки колотушки диски ириски и
      ваще нормально всё
    </p>

    <h2 class="mt-5 mb-4">Наша команда</h2>

    <div class="row">
      <div class="col-md-4 text-center">
        <div class="mb-2">
          <a href="mailto:director@auto-orgul.ru" class="text-decoration-none">
            <i class="bi bi-envelope me-1"></i>
            director@auto-orgul.ru
          </a>
        </div>

        <RouterLink class="text-decoration-none text-dark" to="/employees/director">
          <div class="card shadow rounded my-3 h-100">
            <div class="card-body text-center">
              <img src="/assets/img/me.jpg" class="avatar" alt="Директор" />

              <div class="avatar-post">
                <i class="bi bi-person-badge me-1"></i>
                Генеральный директор
              </div>

              <div class="avatar-text">Ответственный за стратегию развития компании</div>
            </div>
          </div>
        </RouterLink>
      </div>

      <div class="col-md-4 text-center">
        <div class="mb-2">
          <a href="mailto:manager@auto-orgul.ru" class="text-decoration-none">
            <i class="bi bi-envelope me-1"></i>
            manager@auto-orgul.ru
          </a>
        </div>

        <RouterLink class="text-decoration-none text-dark" to="/employees/manager">
          <div class="card shadow rounded my-3 h-100">
            <div class="card-body text-center">
              <img src="/assets/img/jezza.jpg" class="avatar" alt="Менеджер" />

              <div class="avatar-post">
                <i class="bi bi-gear-fill me-1"></i>
                Менеджер по продажам
              </div>

              <div class="avatar-text">Помогает с подбором автозапчастей</div>
            </div>
          </div>
        </RouterLink>
      </div>

      <div class="col-md-4 text-center">
        <div class="mb-2">
          <a href="mailto:support@auto-orgul.ru" class="text-decoration-none">
            <i class="bi bi-envelope me-1"></i>
            support@auto-orgul.ru
          </a>
        </div>

        <RouterLink class="text-decoration-none text-dark" to="/employees/support">
          <div class="card shadow rounded my-3 h-100">
            <div class="card-body text-center">
              <img src="/assets/img/kimi.jpg" class="avatar" alt="Консультант" />

              <div class="avatar-post">
                <i class="bi bi-wrench-adjustable me-1"></i>
                Технический консультант
              </div>

              <div class="avatar-text">Консультирует по установке и совместимости деталей</div>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </main>
</template>
