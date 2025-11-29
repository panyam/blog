---
title: Testing Co-Located Assets
date: 2025-11-29T03:00:00PM
draft: true
summary: A test post to verify co-located asset handling
tags:
  - testing
template: BasePageMD.html
---

# Testing Co-Located Assets

This post tests co-located asset handling.

Here's an image using relative path:

![Test Diagram](./diagram.png)

And here's the same using AssetURL template function:

<img src="{{ AssetURL "diagram.png" }}" alt="Test Diagram via AssetURL" />

