# pccgo

Example of parallel build system

## 背景

C/C++のプロジェクトでは、

  1. それぞれのソースをコンパイル／アセンブル
  2. 全てのobjectをリンク

という基本的な流れがあります。
このうち、1については順不同で並列実行可能なため、`make --jobs`等がよく使われると思います。

が、組込分野の開発環境等、並列Buildができない環境があるため、並列Buildシステムを試作しました。
あるプロジェクトでは、7threadsで、100sec → 30sec程度の高速化は実現できました。

## Usage

最初に、`go build`でServer Programを作成します。
その後、`go build ./cmd/pccc`でClient Programを作成します。

    $ go build
    $ go build ./cmd/pccc

`pccc.exe`をリネームし、`*.join.exe`と`*.wo_join.exe`を作成します。
`AAA.wo_join.exe`は、Server側に`AAA.exe`で処理するJobを投入して、正常終了で抜けます。
`AAA.join.exe`は、Server側の処理が全て終わるまで待ってから、`AAA.exe`で処理を行います。

    # example) LPCXpresso 8.2.2

    $ copy pccc.exe "arm-none-eabi-c++.join.exe"
    $ copy pccc.exe "arm-none-eabi-c++.wo_join.exe"
    $ copy pccc.exe "arm-none-eabi-gcc.wo_join.exe"

    # 上記で作成したexeは、`path\to\nxp\LPCXpresso_8.2.2_650\lpcxpresso\tools\bin`にコピーします
    # LPCXpressoプロジェクトの設定で、それぞれ以下を設定します
    #   Linker       : arm-none-eabi-c++.exe -> arm-none-eabi-c++.join.exe
    #   C++ Compiler : arm-none-eabi-c++.exe -> arm-none-eabi-c++.wo_join.exe
    #   C Compiler   : arm-none-eabi-gcc.exe -> arm-none-eabi-gcc.wo_join.exe

    # 注意) lpcxpressoを立ち上げる時点で、PATHに`c:\nxp\LPCXpresso_8.2.2_650\lpcxpresso\tools\bin`を追加しておくこと
