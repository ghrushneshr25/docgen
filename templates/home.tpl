---
slug: /
title: Home
hide_title: true
displayed_sidebar: null
description: Companion notes and Go implementations for Data Structures and Algorithms Made Easy by Narasimha Karumanchi.
---

<div class="dsa-home dsa-home--book">

  <header class="dsa-home__book-hero">
    <div class="dsa-home__book-preface">This is an implementation of</div>
    <div class="dsa-home__book-title" role="group" aria-label="Data Structures and Algorithms Made Easy">
      <span class="dsa-home__book-line">Data Structures</span>
      <span class="dsa-home__book-and">And</span>
      <span class="dsa-home__book-line">Algorithms</span>
      <span class="dsa-home__book-made">Made Easy</span>
    </div>
    <div class="dsa-home__book-by">By <cite>Narasimha Karumanchi</cite></div>
    <div class="dsa-home__book-lead">
      Go notes and exercises organized by topic. Open a topic below to browse concepts and problems side by side.
    </div>
  </header>

  <section class="dsa-home__topics-block" aria-labelledby="topics-heading">
    <h2 id="topics-heading" class="dsa-home__topics-heading">Topics</h2>
    <div class="dsa-home__grid">
{{- range . }}
      <a class="dsa-home__card" href="./{{.Slug}}">
        <div class="dsa-home__card-top">
          <span class="dsa-home__card-title">{{.Name}}</span>
          <span class="dsa-home__card-arrow" aria-hidden="true">→</span>
        </div>
        <div class="dsa-home__card-stats">
          <span class="dsa-home__pill">{{ len .Concepts }} concept{{ if ne (len .Concepts) 1 }}s{{ end }}</span>
          <span class="dsa-home__pill dsa-home__pill--accent">{{ len .Problems }} problem{{ if ne (len .Problems) 1 }}s{{ end }}</span>
        </div>
      </a>
{{- end }}
    </div>
  </section>

</div>
