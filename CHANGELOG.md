## [1.12.2](https://github.com/Necrom4/sbb-tui/compare/v1.12.1...v1.12.2) (2026-04-02)


### Bug Fixes

* **ui:** ensure timetable uses Swiss time ([#47](https://github.com/Necrom4/sbb-tui/issues/47)) ([a9d68f7](https://github.com/Necrom4/sbb-tui/commit/a9d68f7703934835f95e6abc02e189ec4eb39a60)), closes [#43](https://github.com/Necrom4/sbb-tui/issues/43)

## [1.12.1](https://github.com/Necrom4/sbb-tui/compare/v1.12.0...v1.12.1) (2026-03-31)


### Bug Fixes

* **keys:** disable autocompletion if cursor is not at string's end ([d228ac1](https://github.com/Necrom4/sbb-tui/commit/d228ac1540cb74ed9db94a4473eb1defb21b2791)), closes [#33](https://github.com/Necrom4/sbb-tui/issues/33)

# [1.12.0](https://github.com/Necrom4/sbb-tui/compare/v1.11.6...v1.12.0) (2026-03-31)


### Bug Fixes

* **config:** normalize `ErrorText` yaml tag to camelCase ([899f52e](https://github.com/Necrom4/sbb-tui/commit/899f52e9d85538da47bf188f4c04de37d3db3204))


### Features

* **config:** add nerdFont setting ([f9522d1](https://github.com/Necrom4/sbb-tui/commit/f9522d12f70a6d42348e428fae34d5cf61ad125f))

## [1.11.6](https://github.com/Necrom4/sbb-tui/compare/v1.11.5...v1.11.6) (2026-03-30)


### Bug Fixes

* **ui:** use MacOS displayable Unicode icon for prompt ([4a3eea2](https://github.com/Necrom4/sbb-tui/commit/4a3eea2b2c24658fadbcf47cd1da41b33eba1f4c))

## [1.11.5](https://github.com/Necrom4/sbb-tui/compare/v1.11.4...v1.11.5) (2026-03-30)


### Bug Fixes

* **ui:** add missing `'` from delay number ([02f1bd6](https://github.com/Necrom4/sbb-tui/commit/02f1bd64512157f73977b27ea4f6dda0e67af171))
* **ui:** remove arrGap ([117a3c3](https://github.com/Necrom4/sbb-tui/commit/117a3c381b3d98ea0baa13e62a9eb5e7eae41ddd))
* **ui:** remove extra top and bottom lines around walk section ([b386e38](https://github.com/Necrom4/sbb-tui/commit/b386e38af9a4fdf618d07f25b50d420f7bcb056c))
* **ui:** replace unicode icons by better (prettier and centered) icons ([a2feffd](https://github.com/Necrom4/sbb-tui/commit/a2feffd6ad51fe8fd1d5c2016a4d6f0ec32cb036))

## [1.11.4](https://github.com/Necrom4/sbb-tui/compare/v1.11.3...v1.11.4) (2026-03-30)


### Bug Fixes

* **ui.autosuggestions:** suggest even when not exact text ([d4a787d](https://github.com/Necrom4/sbb-tui/commit/d4a787d36c70c243e2b1a7b6e9d61b681551099f)), closes [#30](https://github.com/Necrom4/sbb-tui/issues/30)

## [1.11.3](https://github.com/Necrom4/sbb-tui/compare/v1.11.2...v1.11.3) (2026-03-29)


### Performance Improvements

* **mise:** improve missing tools check ([b93f600](https://github.com/Necrom4/sbb-tui/commit/b93f600f3ecb859e2dadec76d1b983829f665824))
* **mise:** set text format as variables ([3f47acd](https://github.com/Necrom4/sbb-tui/commit/3f47acda9e44f784701a527089357ce9e2dc3311))

## [1.11.2](https://github.com/Necrom4/sbb-tui/compare/v1.11.1...v1.11.2) (2026-03-29)


### Bug Fixes

* normalize error handling to Bubbletea/Go conventions ([eb65ea6](https://github.com/Necrom4/sbb-tui/commit/eb65ea65fdbaa7062846182bd31ea78562098881))

## [1.11.1](https://github.com/Necrom4/sbb-tui/compare/v1.11.0...v1.11.1) (2026-03-29)


### Bug Fixes

* excessive line when delay is renderd below arrival time ([1141900](https://github.com/Necrom4/sbb-tui/commit/114190035a097f503b1460517cb22d61f38a8f1c))

# [1.11.0](https://github.com/Necrom4/sbb-tui/compare/v1.10.0...v1.11.0) (2026-03-29)


### Features

* **icons:** add differntiation between platform and stop ([c174fc2](https://github.com/Necrom4/sbb-tui/commit/c174fc28c9e14582bb16d6a9a28360c234bf1b13))
* **ui:** add platform alignment logic ([25ea4fe](https://github.com/Necrom4/sbb-tui/commit/25ea4fe48faa408ef2ac731b1aa3010da2b1aeab))

# [1.10.0](https://github.com/Necrom4/sbb-tui/compare/v1.9.0...v1.10.0) (2026-03-29)


### Features

* add update checker ([c52f668](https://github.com/Necrom4/sbb-tui/commit/c52f66841ac1a4c499c6ffc52e3b561895dc2354))

# [1.9.0](https://github.com/Necrom4/sbb-tui/compare/v1.8.3...v1.9.0) (2026-03-29)


### Features

* added THEMES.md file with examples ([c2280f3](https://github.com/Necrom4/sbb-tui/commit/c2280f3b3a71568c16e7a8df56e092c332760130))

## [1.8.3](https://github.com/Necrom4/sbb-tui/compare/v1.8.2...v1.8.3) (2026-03-28)


### Bug Fixes

* **ui.style:** use style of other error messages ([74bab81](https://github.com/Necrom4/sbb-tui/commit/74bab810a2ec1b5190772018e391350d3dc57bc1))

## [1.8.2](https://github.com/Necrom4/sbb-tui/compare/v1.8.1...v1.8.2) (2026-03-28)


### Bug Fixes

* change string according to diagnostics ([ea98295](https://github.com/Necrom4/sbb-tui/commit/ea982954bbb50689674cab768513332e77b6c6e4))

## [1.8.1](https://github.com/Necrom4/sbb-tui/compare/v1.8.0...v1.8.1) (2026-03-28)


### Bug Fixes

* **views:** change icons and text to be more similar to SBB app ([780bdcb](https://github.com/Necrom4/sbb-tui/commit/780bdcb8362d48668048bb45714f678855cd9e96))

# [1.8.0](https://github.com/Necrom4/sbb-tui/compare/v1.7.2...v1.8.0) (2026-03-28)


### Bug Fixes

* handle simple connections without vehicles ([90b31d3](https://github.com/Necrom4/sbb-tui/commit/90b31d3dd2cbc0723677c1043da549cff155f081))
* prefer platform info from the boarded section ([495b7dd](https://github.com/Necrom4/sbb-tui/commit/495b7dda746770425392735a34728d355b2b5391))
* use boarded sections for simple connection times ([aaf1172](https://github.com/Necrom4/sbb-tui/commit/aaf1172e8afb4d35cc5fc67ecd2220ecb867f7f0))


### Features

* move walk time into the simple connection timeline ([59c1eda](https://github.com/Necrom4/sbb-tui/commit/59c1edadd3bc248f616b92493dca8d2766313f41))

## [1.7.2](https://github.com/Necrom4/sbb-tui/compare/v1.7.1...v1.7.2) (2026-03-28)


### Bug Fixes

* **views.docs:** document scroll key shortcuts ([2ea4bfe](https://github.com/Necrom4/sbb-tui/commit/2ea4bfe54baefd3290605c6e710e719468e209b2)), closes [#15](https://github.com/Necrom4/sbb-tui/issues/15)

## [1.7.1](https://github.com/Necrom4/sbb-tui/compare/v1.7.0...v1.7.1) (2026-03-28)


### Bug Fixes

* **config:** add missing colors ([d86e9a3](https://github.com/Necrom4/sbb-tui/commit/d86e9a337b5d3f6a540b9a854446c5565913add6))

# [1.7.0](https://github.com/Necrom4/sbb-tui/compare/v1.6.0...v1.7.0) (2026-03-28)


### Bug Fixes

* Change warning color from dark green to red ([d69dfd4](https://github.com/Necrom4/sbb-tui/commit/d69dfd44cf78e89109422fdff14af8d54f408a82))
* **config:** prefer `$HOME/.config/` ([4cd361c](https://github.com/Necrom4/sbb-tui/commit/4cd361c884418702512ced1ff7b201652c228343))
* **config:** revert back to original colors ([8f54d5d](https://github.com/Necrom4/sbb-tui/commit/8f54d5d49471266a0f4e6b39306b18f53420f159))
* use new color system on newer master branch commits ([8074c14](https://github.com/Necrom4/sbb-tui/commit/8074c14ca289cb7b4900a3c82d8c2da82f1aefa3))
* **views:** replace old noStyles (now `lipgloss.NewStyle`) by Text color ([a5421d6](https://github.com/Necrom4/sbb-tui/commit/a5421d6ce3c171600f22c75e9c40672f1309689c))
* **views:** revert back to old color attributions and add new colors ([bc5fe8f](https://github.com/Necrom4/sbb-tui/commit/bc5fe8faf13ec110be062a83345c971586bdbe7b))


### Features

* load optional theme config file from ~/.config/sbb-tui/config.yaml ([4c67797](https://github.com/Necrom4/sbb-tui/commit/4c6779792249a35d4226ff10d8518aaeb6179369))


### Performance Improvements

* **views:** define styles the idiomatic way ([1d74efe](https://github.com/Necrom4/sbb-tui/commit/1d74efe2d8030add6ccea99cc14675ca92fed96c))

# [1.6.0](https://github.com/Necrom4/sbb-tui/compare/v1.5.1...v1.6.0) (2026-03-27)


### Features

* **views:** add AcceptSuggestion keybinding for date/time fill ([15948a9](https://github.com/Necrom4/sbb-tui/commit/15948a969e721adddc5f426ecb557742657576e9))
* **views:** fill date/time inputs when incomplete + ghost completion ([9e28b11](https://github.com/Necrom4/sbb-tui/commit/9e28b11db038a22207b508f8c99a52cafd51c551)), closes [#14](https://github.com/Necrom4/sbb-tui/issues/14)

## [1.5.1](https://github.com/Necrom4/sbb-tui/compare/v1.5.0...v1.5.1) (2026-03-26)


### Bug Fixes

* **api:** debounce autosuggestion API call ([eceaf0f](https://github.com/Necrom4/sbb-tui/commit/eceaf0f6211c0cf14b2f702d9f68883e0f0a218c))

# [1.5.0](https://github.com/Necrom4/sbb-tui/compare/v1.4.0...v1.5.0) (2026-03-26)


### Features

* **views:** change to Swiss date format ([ea7c466](https://github.com/Necrom4/sbb-tui/commit/ea7c466950524e9bb6de08a9ea3d765b36bf70d3)), closes [#9](https://github.com/Necrom4/sbb-tui/issues/9)
* **views:** prefill date/time fileds ([1cf3173](https://github.com/Necrom4/sbb-tui/commit/1cf3173da99c3ed1fb40531c2b964f9b320b3687)), closes [#6](https://github.com/Necrom4/sbb-tui/issues/6)

# [1.4.0](https://github.com/Necrom4/sbb-tui/compare/v1.3.4...v1.4.0) (2026-03-26)


### Features

* add Windows .exe release build ([bcee98c](https://github.com/Necrom4/sbb-tui/commit/bcee98cb31386abe21d082a2850a03ede490e47b)), closes [#4](https://github.com/Necrom4/sbb-tui/issues/4)

## [1.3.4](https://github.com/Necrom4/sbb-tui/compare/v1.3.3...v1.3.4) (2026-03-25)


### Bug Fixes

* **views:** date/time string moving left when typing last character ([1040de3](https://github.com/Necrom4/sbb-tui/commit/1040de3085082bd9807d09dfabb874ee12589b08))

## [1.3.3](https://github.com/Necrom4/sbb-tui/compare/v1.3.2...v1.3.3) (2026-03-25)


### Bug Fixes

* **views:** clip text to prevent line wrapping ([2294cab](https://github.com/Necrom4/sbb-tui/commit/2294cab918eaa337e95235a136573e6cbc0dd158))
* **views:** enforce fixed minimum terminal size ([76de3e6](https://github.com/Necrom4/sbb-tui/commit/76de3e6b02125ca9d2cc31c0bc7cabcf8c467b94))
* **views:** fixed size detailedResult window ([213f9ee](https://github.com/Necrom4/sbb-tui/commit/213f9ee2125ef88ebb2789bc57bc842e4a74b574))
* **views:** reduce date/time inputs width ([4ec1b5d](https://github.com/Necrom4/sbb-tui/commit/4ec1b5d5cedf024778a12170ba87a9c00d163075))
* **views:** revise shortcut bar names ([e91d021](https://github.com/Necrom4/sbb-tui/commit/e91d0217547811db931e9dc296ac4275c20b7c25))

## [1.3.2](https://github.com/Necrom4/sbb-tui/compare/v1.3.1...v1.3.2) (2026-03-22)


### Bug Fixes

* **views:** let header take whole window width ([40ce6b9](https://github.com/Necrom4/sbb-tui/commit/40ce6b95088d7533b641503afac669d0618aede4))

## [1.3.1](https://github.com/Necrom4/sbb-tui/compare/v1.3.0...v1.3.1) (2026-03-22)


### Bug Fixes

* backspace removes separator together with its digit ([613c2a8](https://github.com/Necrom4/sbb-tui/commit/613c2a82b9f76db169690f39260dc01dc385bb68))
* clear stale results when searching with empty inputs ([e431e25](https://github.com/Necrom4/sbb-tui/commit/e431e25b5d4a61da2905c526b27abadf19b8fdd9))
* correct date and time input validation ([083ffe1](https://github.com/Necrom4/sbb-tui/commit/083ffe116f5a1d02a9c63aa17131abdd036b7298))
* prevent crash on small terminal windows ([3d0378f](https://github.com/Necrom4/sbb-tui/commit/3d0378ff6fdb8dd91598909a15e8acd7cf51039e))
* show 'terminal too small' message instead of rendering broken UI ([39e1ac4](https://github.com/Necrom4/sbb-tui/commit/39e1ac4cbf464afe2d62a1bb08bfa95a6ca88921))

# [1.3.0](https://github.com/Necrom4/sbb-tui/compare/v1.2.1...v1.3.0) (2026-03-22)


### Features

* **views:** implement sbb-logo on startup screen ([27f442e](https://github.com/Necrom4/sbb-tui/commit/27f442e942c3f5d3fb409df893e00635ee2ddaab))

## [1.2.1](https://github.com/Necrom4/sbb-tui/compare/v1.2.0...v1.2.1) (2026-03-22)


### Bug Fixes

* prevent index out of range panic in date input ([4a67fa5](https://github.com/Necrom4/sbb-tui/commit/4a67fa52f3ace4f10495ca05023ae34a3443163c)), closes [#1](https://github.com/Necrom4/sbb-tui/issues/1)
* prevent index out of range panic in date input ([a220f6f](https://github.com/Necrom4/sbb-tui/commit/a220f6fe9267e8650d8e1d7e1bd5d6fac819d6ac))

# [1.2.0](https://github.com/Necrom4/sbb-tui/compare/v1.1.0...v1.2.0) (2026-03-09)


### Features

* add base flags ([dd19f14](https://github.com/Necrom4/sbb-tui/commit/dd19f14177d64bd07a0b4118e9407fffbf37b1df))
* **ui:** add `no-nerdfont` flag and create two groups of icons ([4eeee6e](https://github.com/Necrom4/sbb-tui/commit/4eeee6e375ab958e7b39938cc8e5d105f721183b))

# [1.1.0](https://github.com/Necrom4/sbb-tui/compare/v1.0.3...v1.1.0) (2026-03-01)


### Features

* add station suggestion ([24257b3](https://github.com/Necrom4/sbb-tui/commit/24257b3a1d08af6413528385a406bc57ea10c3a2))
* **view:** add help bar ([21dd62f](https://github.com/Necrom4/sbb-tui/commit/21dd62fb9a7032663ba5322d869b0ff783055d02))

## [1.0.3](https://github.com/Necrom4/sbb-tui/compare/v1.0.2...v1.0.3) (2026-03-01)


### Bug Fixes

* use correct homebrew-tap master branch ([95e1990](https://github.com/Necrom4/sbb-tui/commit/95e199054433b75269f256ae962c6d5b84fe2705))

## [1.0.2](https://github.com/Necrom4/sbb-tui/compare/v1.0.1...v1.0.2) (2026-03-01)


### Bug Fixes

* revert goreleaser homebrew key to brews ([47b2da3](https://github.com/Necrom4/sbb-tui/commit/47b2da32c5f97ef04018f74d6fda8de4cf17a63a))

## [1.0.1](https://github.com/Necrom4/sbb-tui/compare/v1.0.0...v1.0.1) (2026-03-01)


### Bug Fixes

* initialize homebrew tap ([02ac0be](https://github.com/Necrom4/sbb-tui/commit/02ac0be682a264ac5ef8b93dfd555174a4abba28))

# 1.0.0 (2026-03-01)


### Bug Fixes

* **api:** remove `date` and `time` when not specified ([559ea12](https://github.com/Necrom4/sbb-tui/commit/559ea12226b740c8ffc4fe9ba6586f16518efcb2))
* **models:** reformat SBB date to time.Time data type ([05975be](https://github.com/Necrom4/sbb-tui/commit/05975be35dcc6877621e234f6272bece7ce42a56))
* **views:** use variables for size handling ([1284f23](https://github.com/Necrom4/sbb-tui/commit/1284f23b8116cbec43bfb23288afa7fbe0f6172a))


### Features

* add api package ([47972cd](https://github.com/Necrom4/sbb-tui/commit/47972cdc6291722161076d6d3a606ff1946a6ac7))
* add delay information for beginning departure and end arrival ([cd4abf9](https://github.com/Necrom4/sbb-tui/commit/cd4abf95e8b649a7762de89d8d83f80fe44d3e54))
* add homebrew-tap release files ([04c758d](https://github.com/Necrom4/sbb-tui/commit/04c758dbb296671ca81cc0669ecae3ac8ca57f2d))
* add mise.toml with latest go version ([46ccba8](https://github.com/Necrom4/sbb-tui/commit/46ccba877166017d844aeb8e3e800f4a587b6b08))
* add walk time ([665b6e6](https://github.com/Necrom4/sbb-tui/commit/665b6e60d154110c6c035ef7d3b6f3c2e508bacc))
* **api:** use the current date as a default ([6aa0cd5](https://github.com/Necrom4/sbb-tui/commit/6aa0cd5d7c94f82a6bf6c9e12eadfcc0bbb6c1a8))
* **connections:** add date and time options ([65cc412](https://github.com/Necrom4/sbb-tui/commit/65cc4128cfcbe6c25aeead3a458b4dfb41633635))
* **connections:** add isArrivalTime switch ([9bb2310](https://github.com/Necrom4/sbb-tui/commit/9bb231041c69966c52825cffea684f81824c9785))
* create first structures ([5d4ade5](https://github.com/Necrom4/sbb-tui/commit/5d4ade578b179b9e192936625342cc2bd05c9ad2))
* implement walk and fix code accordingly ([69efdef](https://github.com/Necrom4/sbb-tui/commit/69efdef3d9a03611727f45526a9833ef4a51b83f))
* introduce views using Bubbles textinput example ([725b1ef](https://github.com/Necrom4/sbb-tui/commit/725b1ef5ef11a30007ae255f31c07adf7654992b))
* **model:** add capacity and section.journey ([5968111](https://github.com/Necrom4/sbb-tui/commit/5968111b8e60b74edf86b8ba8014588fffd3ef9d))
* **models:** improve structures ([a88cd38](https://github.com/Necrom4/sbb-tui/commit/a88cd381b1e89115cd5303b6bef87889a524ca9f))
* publish to go and homebrew ([3323c67](https://github.com/Necrom4/sbb-tui/commit/3323c67ffa32b544a9b7edfdc8f01112de5e8178))
* use Claude to add fullConnection render ([75cfdb3](https://github.com/Necrom4/sbb-tui/commit/75cfdb322b4141a5626c0bec03600ea4a249af0a))
* **views:** add input restriction to time field ([3323c82](https://github.com/Necrom4/sbb-tui/commit/3323c823ffc4e4827a8ae0d90d93a4f082c5527b))
* **views:** add renderFullConnection border ([799af5a](https://github.com/Necrom4/sbb-tui/commit/799af5a991b242101deb240c8cf8f1501cd2621e))
* **views:** add search button ([2d4ba6f](https://github.com/Necrom4/sbb-tui/commit/2d4ba6f34ad765cb3955aa5fa9406ed0e7bd5f4d))
* **views:** add second result window ([386c2a2](https://github.com/Necrom4/sbb-tui/commit/386c2a2585455f26aaac836df1a51e35384fecbe))
* **views:** add Stops visual representation, fake vehicle icon, company icon and model icon ([a300824](https://github.com/Necrom4/sbb-tui/commit/a300824080e55fc303b8110318be4d044522bc60))
* **views:** add swap inputs button ([3b25e43](https://github.com/Necrom4/sbb-tui/commit/3b25e4385379695ea715f1de3b319b8bee7a82cf))
* **views:** add url link under walk duration ([738765d](https://github.com/Necrom4/sbb-tui/commit/738765de1e81c98c7d9a5054166212c61ed5423c))
* **views:** dynamic length for stopsLine ([bca8d87](https://github.com/Necrom4/sbb-tui/commit/bca8d87467003658237b16fd1f2cd62f7ada1617))
* **views:** dynamic number of displayable results ([1a2e978](https://github.com/Necrom4/sbb-tui/commit/1a2e978533afbaad4c6d70a3173e5d44675f8cc1))
* **views:** improve UI style ([fc64da3](https://github.com/Necrom4/sbb-tui/commit/fc64da3f321e55d07ab135fc74dbd8f1befc6710))
* **views:** modify example view to a simple travel searcher using API ([d6bb91c](https://github.com/Necrom4/sbb-tui/commit/d6bb91ccc71bdbcff4d55f6f167586bba3ed9950))
* **views:** put results into boxes ([fc01c99](https://github.com/Necrom4/sbb-tui/commit/fc01c99d86810d0a14029431998f0d590335eafc))
* **views:** re-add title ([c846841](https://github.com/Necrom4/sbb-tui/commit/c84684159f54e0656bb2966293ce102fd64da006))
* **views:** start implementing full screen ([31ce207](https://github.com/Necrom4/sbb-tui/commit/31ce207d43ed86962a545e0197a6829756eaf7bb))
