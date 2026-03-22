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
