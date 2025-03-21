# emptycheckパッケージ

##  構造体の各フィールドのemptyチェックをする
* デフォルト値である場合にemptyと判定する
* ポインタの場合はnilの場合にemptyと判定する
* 構造体のフィールドにrequire:"noRequired"のタグを設定するとチェックをスキップする
* isEmptyメソッドを定義することで独自の方に独自のempty判定をさせることができる