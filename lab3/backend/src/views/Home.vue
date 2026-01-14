<template>
  <main class="container">
    <h1>Авто-Оргул</h1>

    <p>
      Тазы валят <br />
      Шеснарь дует двинарь газует <br />
    </p>

    <div id="mainCarousel" class="carousel slide mt-4" data-bs-ride="carousel">
      <div class="carousel-indicators">
        <button
          type="button"
          data-bs-target="#mainCarousel"
          data-bs-slide-to="0"
          class="active"
        ></button>
        <button type="button" data-bs-target="#mainCarousel" data-bs-slide-to="1"></button>
      </div>

      <div class="carousel-inner rounded shadow-sm">
        <div class="carousel-item active">
          <img src="/assets/img/2110.jpg" class="d-block w-100" alt="ВАЗ 2110" />
        </div>
        <div class="carousel-item">
          <img src="/assets/img/2112.jpg" class="d-block w-100" alt="ВАЗ 2112" />
        </div>
      </div>

      <button
        class="carousel-control-prev"
        type="button"
        data-bs-target="#mainCarousel"
        data-bs-slide="prev"
      >
        <span class="carousel-control-prev-icon"></span>
      </button>
      <button
        class="carousel-control-next"
        type="button"
        data-bs-target="#mainCarousel"
        data-bs-slide="next"
      >
        <span class="carousel-control-next-icon"></span>
      </button>
    </div>
    <section class="mt-5">
      <h2 class="mb-4">Отзывы о нас</h2>

      <div class="row">
        <div class="col-md-4 mb-3">
          <div class="card shadow-sm h-100">
            <div class="card-body">
              <p class="card-text">Машины нет, магаз норм</p>
            </div>
            <div class="card-footer bg-white">
              <strong>Лёха</strong><br />
              <small class="text-muted">г. Пермь</small>
            </div>
          </div>
        </div>

        <div class="col-md-4 mb-3">
          <div class="card shadow-sm h-100">
            <div class="card-body">
              <p class="card-text">
                На приору ваще топ, жаль гнилая. Купил 10 литров фриза, до дома доехать хватит
              </p>
            </div>
            <div class="card-footer bg-white">
              <strong>Диман</strong><br />
              <small class="text-muted">г. Пермь</small>
            </div>
          </div>
        </div>

        <div class="col-md-4 mb-3">
          <div class="card shadow-sm h-100">
            <div class="card-body">
              <p class="card-text">У меня авео, побили(</p>
            </div>
            <div class="card-footer bg-white">
              <strong>Тимофей</strong><br />
              <small class="text-muted">г. Пермь</small>
            </div>
          </div>
        </div>
      </div>
    </section>
  </main>
</template>
