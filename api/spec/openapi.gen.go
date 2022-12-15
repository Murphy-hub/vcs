// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	externalRef0 "github.com/trustbloc/vcs/pkg/restapi/v1/common"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+x963LbOLLwq6D0fVVJqmQ5O5fdsz5/1mNldrSbxF5b8dSpSUoFky0JYwrgAKBlTcrv",
	"fgo3EiQBkrKtzMzZ/RVHxLW70Td0Nz6PErbJGQUqxejk80gka9hg/edpkoAQc3YL9BJEzqgA9XMKIuEk",
	"l4TR0cnoHUshQ0vGkWmOdHvkOkxG41HOWQ5cEtCjYt1sIVWz9nDzNSDTAukWiAhRQIpudkiqT4VcM05+",
	"xao5EsDvgKsp5C6H0clISE7oavQwHiULymgSWO+VboISRiUmVP2JkW6KJEM3gAoBqfoz4YAlIIxyztgS",
	"sSXKmRAghJqYLdEt7NAGS+AEZ2i7Boo4/FKAkGbIhEMKVBKcdS1vAfc54SAWJACKGZWwAo5SoEyPqgCQ",
	"kSVIsgFE1PYTRlOhVqM+2TG9+YgZQU3YNdG8e1wfHeHBOSw5iHUXTm0TM8oYbdckWaMEUx/k7EahBFHY",
	"1uYUQQiKhOUB9J5fzGfn70/fjhFZIqJRkOBMja62ojs5RFVUlWQEqPxvxOQa+JYIGKPLN//6MLt8Mw3O",
	"rZe1MD+HNqu+OOj5VBwYTEPvl4JwSEcnP9UPR22iT+ORJDJTfUPnshyY3fwMiRyNR/dHEq+EGpSRNPkm",
	"IaNPD+PRWUmXM7UkrjZQP58pEXmGd+pPImGjf6uP/VBOhjnHO/X/jCU4A69pBSyKN6EPD9V+WmuKbEbt",
	"hZgW9a1cGoR2Mafz2fQMVT0cDbTZU0rS9jjT2VSRkKFay5GqE47WWFjuccMKmgZpZsn4BgeW+L3+vTzD",
	"1aA3oM58lHD0utlSDfj/OSxHJ6P/d1xx8mPLxo//8eP8QrcrkdYiWI9Yh0/foFv9dayB9ymEWIeh/cl0",
	"iORpI7dL+OQSK9YbYVan6B9X5++RCIgJc5RFcSPUbqjMdk3Whb1VTNC7D1dzRRQ5BwFUGinhgZgIRJlE",
	"HGTBaQTJUTkWXeUBhNnZ04WZXi55TolWAVLNxiicL0cnP7UZUIt7feo4jT5Ua6tc1o6pZQGdcGmcEDtj",
	"bd2Ro/Joln4lsSxEe1fe0RC6SftgiLJrW9pZttG9PzuAbR7c2VWtSXBfQfZu+p3nAXyd6z+E5gKqrz4N",
	"NazUtzlsL31bUEsZuIspSc8YXZJVWKiYbx0M62/qbMN9YOv2Q/BEZoTeQrpISRqghgvDjowKTSj6eSte",
	"mq6vEOPoZ8Folr4023plCV9hrdQG9jxrdV3Bk/slbAaRfAp3ODcU/+Y+WWO6glPfGjhjKQxQBMD01Yyt",
	"kGuUsBTQkrONOdQcMfVzCw8sXygKH0A8ZUuPgHoXPJCaOsaJSUn3BW2eCgJ5vzC6Uc/h0c2GbX4A2r3d",
	"/wA4k+uzNSS3e+13rfuhRHWMMsCk4ByonJNNYNAz8xFpOWVlQGV8OmEySrGEI9UmaLBEmLNhKUoZ+DgS",
	"hVbtP46U9WImUB+KHGGaIl5QJeP7ZY2dysNBCHRdUDcg0xDToJ9RIgmWoDStb85mA86Z69FSzpSOrxQw",
	"dBlTwWtG/iIFiUkWkmqFkGxDfgWBtmss0S2hqUKOtR2NLYG2mEqtna/InVaLrs+uwlpMhslmkWKJQ0Rl",
	"gKt3dsHhyNGyUgLU6fk+Y9uJGtps9wr4HUmU6SwFwgKdX+ieW5xlIBHO84wkendtaViuBGiaM0IDQD5T",
	"35H77ujR7lef4+0aeE2b10MitTnPWKnsYLyUwJGlvmWRZTuEE7VlzSN6bXFjPy+IRfmCWBQvCp61l//h",
	"8q2vR2lasF2VRPL3hdGPGmQTNMe3IJQmnag9JYCYOht24i1k2S1l21JtRTnmeAMS+ATNluiGqePfsUh9",
	"vFqDYQ5aQc85uyOp0qSNZmwZjBup2oXa2ZZkmVPIUaJJNNKS0FKrzIGS9Mg1O3LNTo6Pu+BdrnSIl8vQ",
	"3vGaZSlwnwQNxZohUbX5REvmgps2Hy7fhldSkthCwibPNGADBvTcfgwYmoYWrUa9XZMM6oSYMJpkRWoU",
	"cSK0EcFxogaelG4a7e5RA+ecLdUQRJQ7MCZQoURCkUmSZ/Xp7crClL3imMqIp8ceuARTRyEO37qX9gIJ",
	"JNecFau1WbtHlnP1/6qhdyy10WYA4UtqWveLKn5S94ZqMU4oUrvhSEjIhab+NgmnsMRFJtV8dV6rhgjC",
	"wVd/gpR2h7MCrJ1Z+tUaXF/RnWLROf6lAOeSMwccScXBlaizFuuNYuZazBY3R9be1os1Hj29YXfYt0Su",
	"I/OpHSKrKSMBUonStNArzjncEVYID1KVLxApRkPuQCBst6bgXcfhGBFpbHyiKRTU/wl1q3aLPq0v2ko9",
	"t/0AiIT+4CBezWcWYt0K78/nJa0Qimq6lRFJy4xtjS2fczjCpcBaGDoRzi0RxLdjchHSPzN8RVTMUNOw",
	"RaLeBtznoKSfkon2+BmazoErtqdQoDlPnYidXw5NDY3qQ9F0Pfd6gcv16e9i2MJ8l0f7YCn8V1K0vj7D",
	"vye+jRSxoiuPaSGAL3JCF5Xi9kit4zvGMsDU0qnIISHLnWb5a5BrdQicx6LafO7bgFrQqvWgi9l7hDOm",
	"+roz5a5zDNVqH1ednix41FIqDN2YNdXMvYgOOVD5b/XudwsOUUFjdvdAXcauJiLhPX5uOWLFlnIs1PHJ",
	"4E6JAEKNKqDQ0WCMLDC4hja6KvKccSmMfvPDfH6B/v5mrnms/s8lpIRDIid2WoE2eFd6I/91aTDn6QiO",
	"oWo9UUFQEYWmcKGknFYt5RoIRxt2o47Mj6VCG74nuQ8rAzWwOLbnKcXmsDHOIbNeiiWiAGnER+qOUtDV",
	"4VOqAdvfgQLXYup8foFyo4aVsO23rsKkMW4bvjGSfQzFX19MrWVSp1P/IE9hqdfG6CwN8qC84DkTPbcx",
	"oWkDRkqjmX8eO6w5z3AMUMts2u9eCA5nO3+K7iIKe7UTBfKKO0yDFmDFVSwn7fIy1l3TUV/sUN+bmmBP",
	"nxszftG+26HG1q03tQXzsMs6BLehpB2e98nu3WStGBJdhVSCNc4wXWnNB6ep0TKtxcCWMeNGsYnwHW/q",
	"WTNmCKVBsg2RirOInZCwMf4ibRFahtdjRFUO/C6shdzRD+NRyjY4xASn+vc99n0HnCwtL34Hcs0iIPhw",
	"OXMQaHcx/N1ozSEILQkXEkH61bff/umvKC9uMpLoeyi2RNPZFL20ckGrPsamm86mr/qg+RClT0dkA0m0",
	"vDhtMbSftwGfTBkvgK7IikKK/vHjXKny5YWb2lp16Ra/0I1o3NX4+orqKnBFZaZS3SfIOiwzo3Uwmu2Q",
	"MNoCpF5DRRQvft7KF/1Sz1vcWIPA4wQlrIZeWZ0rTfbCGTYixm61EqoAZ1TbHBMufH2lNI2M6VyQLLXe",
	"DMYhbFigl5ffn/35L9/89ZVREQ2R6U7WRjbamTFSnMdOK8f18bTpHuBA1ssUFsD2q4CEQ1gItgyvuMkz",
	"1NZoILI+w9hbcXN9bi4P003EDTxMFxxyzEE7LpWcOI3oBDGZa/sj4/lUIzQs3v19yZbBThSD3TA62eFN",
	"FuS2tYmmdoCGS2Rf+/la07OLDxBGzf84Uvr4x1G3oftMWA/dUg3C0vNgvN94G4DyaJBHDedxP7o5/C9E",
	"4/jXz7nrHsRKfSZeEXKX+G6eIa2YizWki+Bw+2/g4vSye9kxu4xjKoxnFZmAJ2eDASryhG3arhH/4noP",
	"vb0E1TiGrIA5NYyk9qTPDk09QIv/p+LHBlFBuOcBA8nujSRqhmC08fVYRF+CKDK5N7pjzOYgYUgVUlvE",
	"EnHXSr4LYOTywxtElv7tn40524FE+A6TDN9k4JzW1rI9v3B32+aSQmvchKZKvYfqjlMy0wE1Y+oQoUIC",
	"1nexSRuE6OUUlsB5LX5Ku2ZeRRyIEVN0XMVUmf13kYvF+lCiKcQ6JPeHqCqFWDckle0c5xm/iZISC1sZ",
	"R5bjQ7cHPHtAGdL9NQPdbbA20BWVaIM9abG50X5yLBEH61IT9ehEy9icGaFsXy9gEQuElXlHJLkDF+eo",
	"zk+9RxXrKBCWesCUCKVsWz98LMcB3RTSHES5y0mCs2xnbvAzrGZU5t2acYlewmQ1GaMbkFsAir7VzuA/",
	"v37tFvoqFsBvVI2Ck1j4frUJrRQoaJs7WRZYdHkNz4SE1PIRDTIFJ0HoKoOjQui0AOBgA1YNfEUOiYZi",
	"zRvdvlcL3xv1Chh/q7W0iAZ9xwhzqIl7JRl/VICakIzvG5qlmgVtgkedfz2aB47urQw87LFB9ojjegxk",
	"OoLW+ra3n0L5IU+xhKZrLorvzuYl6QvJi0SaqxjVQe3++iwew1YlUkzD/ocneho7xPFsOgqM71FRN4AG",
	"QvkaZ0QNc1FhDNKBB+vO9LV3/60bTMUpc0LbUA1cnwdvmlBjxD3vqj4I4G4Bfa649oI8QPfC6Omw7hfS",
	"jwV2PNDkPNd0D3H3iAgJ+IwIrUg3VmE77HNjH7GXpGcvVabLED5brqFl6j4bCoGT5a46cy4ANGgAmcZB",
	"ld3zuS4xyQoONprWKoehmxhIbkO3MKqX3mYQj8A54+1ub9TPaANC4BU8+s7i2muDNrpR/2EzG3ErC07k",
	"I64D4F04M6NGsNZ3I+lhzF/dH/hesgmB/S4mg/B7NPQHXU7eNc/Ooe8mn+my7yEOtSH3ZZ2AGyImSg5T",
	"8xWIPjpWp0rUvOH7UJN/KLvyU6Ib2hMkfsLNEA5cC876w/DgTr7ZOp0xmDwBtH1ssgbWbgLbi035aygZ",
	"1bgWmfNMSVh7M9y24lgtqRMlj2GZITgMYZr+qvZmm/rT74Bvhjb/BPjtyzv3oO1HMc/Yce1nn8FdDYbM",
	"j5Bl/6RsS89zoLPpmZ+UECIu1ag/qTJ+QdiV4DC0HSkLKwwzrW3Rg9ogYlGGawTrL1RzdF89N8bZw9jR",
	"lSc6ttx5vedd45GypEN9xI6F1qAZ+bkcx2NjXcRSUsMw00ltj9AlM841KnGiYQAbTLLRyWgNWcb+Jnkh",
	"5E3GkkkKdyNX6WI0Vz9/l7EEScAbRX06dHi0ljIXJ8fH9W6KJzXu31z367MrF7VTLxtgo4CVEe2fdVQI",
	"pR38+PUZuj47Or2Y+bHcBibfXOt7YskS5kczHrtD52fXmH42AWs0HmUkAcuS7E5Pc5ys4eiryevWJrfb",
	"7QTrzxPGV8e2rzh+Ozt78/7qjeozkfeGgfj8gugLKc8Ycal0L6/Prl4Z+0sYQL2eqIm1UQEU52R0Mvp6",
	"8lqvJcdyren82E9kPPk8WkEodksXYBDO1R5JF1UMBLvg2tHfQf7gDV1Rs572q9evHeWAOT1eePOxEu5V",
	"qaU+RhFK3dT02eB9/9RnUhSbDea7MuUTndn1hTM7H8ajY0sCHubFsU08qlwmeuVHzveVs5CrzSX6BtMn",
	"mp7aMlyhDdsB2dLWQfodS3fPBujeaR8eHh4OiOj+POkhaH8cEjwCqXhghDZyc7d6pG+Vj1IssaaSX4+8",
	"+JcwgdhbWYF0CEw4hMsP6vPCrWsRLm2SsSNHIpYOQS2DgqUOTDHDImKGUM3QALtH0UnN6xSmjA82GaSM",
	"LvDkXZkQLFl5MVLPH7UpojbRp54lEyOVWiTJIQmkmucLUUMz5mEv/NfiawZjuhDrhqTo5QUtjNsQcD94",
	"Tefn6Wsj5HvZtV5SZ2+eD6SB7UiowqGQ3hMZESeBPgRFw0r2QZSQjO8n0/XlqHiqRO+7QT4EKrrnPPBZ",
	"7LlTHnIkHwP5fWjB3tDBUf1mrIce3M2UiF7rFd49Zp0KBtxqHYIQeqc9MC3036MOIYfhgO8hAlshQRx/",
	"tn/Npg/HnufBtNMU4CVo/BSpfuByc0ykKFFflPFVWcPlJCPfYSB5AWMPfk3nwqdxhP5mzey7EM9nQjaS",
	"fw7F70M5cM9ATg1/zwDy0AtByVDx3UsEVbWi3yMVmFgS4auBMUNdEYNHB2UNukNQQ3eIy29CF52QegYK",
	"Of5s/p1NH7r8K5zAHYhmaHiHcyWEst+QEsfhql16lMAkovq6F7V/YeIYgJi9SaSmY5QVjRhJk98tM/Hq",
	"M5CyPgPxi0fMgv5X31dKqC4sZiN+61e2IlYqKlQEu2yqq45mbFtT9/wqCe1z45Le63cLar5Dib9wSY8D",
	"K1SxugqD5GRfTZAemvdJfbKFLDvS5bKObQmvpHlB1elwrnVqo/Ncf7blQw8Iz85LtmF8xHhlavsJAbKH",
	"h5cnP33Wg6/QpxjQcak1RzFjHJNfTV6H3V+ucL09S6auoK7hVdbmalZ18hPxGuglaVIaBH3irTdjVIPr",
	"lwL4roJXM+nzCfJuXpXGMlX1l8QY+qF5/ZTiJ8x5isrLf5QCJ3eQljVojNlTXlu58mG6so3NFQgmCIxt",
	"op3tmSK8Uhxamqpp0Q2xFBZVJMITd2VCD82at7iqeWb2aKvyuMmGLWlhxhztjdNgsgm3hYOMXqDsyiO8",
	"AloWLTP4fSHKhrW6ja6wWrZDICS+yYjO1ykrSgWntEXZahXYVkRIe3ubc6bPF+OmpNkG37rm0TyQ8Ikw",
	"C7bpH3sCyzwfUX8Wo2dCk7e9H4FQVyLPlD3w60dZ2EiGNpiYGpOmSpzL+PFzlHRRS5xlNzi5NZpJEPS2",
	"ep0w5e3MnLYGvMWuhbRHCGrIOjWYCapidVc/nH94Oy01GxuAdadYhy7iwoQ4EkRWq10yvgK+iwKyDJp+",
	"PH27XDalmN3BzpC3+w3fsEI2FGHTwhaBKCvJmmL/E/TOlZOMTOIpdob4dYKslpGL+gVGibEafghFCTbx",
	"PYHKlSIGqXD63l6QM5f+LwSqImooJNKVy/hw+dag29XYJVmmC+m5hDZ2B3xXHlrN2iTwDaHgAfSFAlGO",
	"b0hGJAGhydUxETFBl2/Ozt+9e/N++maqIDHdUbwhiS9aL7uPnpllUQagPOoIaqfaWt9FVJTw7vR/9HaJ",
	"/xpEedRsSUBJNuRXKA/OC6GL73ECNIFn2J1OJFmbSJq97FavXKeV5Dv7ZBFwzVAs2lz9WLiXLlexYQQB",
	"n6DTaHlMJY6rZMUcC1uqEtNg2d+SDTgBX5liFeRtJmGryq9fOVQX1lNdqhKaZok1ntXeybyac1MIiSS+",
	"1SYeU5yeFdTWKC0HtZnUqwIrBRDsOyCcrAhVn+0+iLCDjlHCiixVHAFThKVUTDmCWy9H5PHug69ff9Vh",
	"IdwfbbfboyXjm6OCZ0CVOpHWTYZwimGsLFFbvGg9ZlUW2+t6wSvWW+u7JkfTJLhmO1urmmh1z2bHK3FI",
	"JFk5s5kTcau4Zgb4NvJ6VDjHyG3HVRP+aBp+HHmktsVl+UynaVqpHKlkqvYG9ziRlg5taVlflzUStD+m",
	"2mV69Xl6vmcFTRtGmjZ4+y71qwzW0mgacn2v5YCoCU5CXfViwxwU0dcmL4u6tq2ig9/Ntx8mOrDLIpBE",
	"MMSwbjgpuhGVYx7H0JkhUQE0deE34ZRto/Jlu1ZxaKcuKjG9AimaqfBVKWDFJn3lB4t2nrdL6vbkJ68q",
	"I8dL8LSJJZisvd994t7McOALBf8GCmj0LYBIzcCgD6E9SN3ePvl9eAZ6luls8JNnsPgfW4D9Pxrdb6/R",
	"BSLyPafDyb+ZF+YLlsbb22EzVC38j0cmXMdgHaxK9zsznltLr/sFTv7wvo++ejcd5T7rYjZkWbSV4j89",
	"a7BkrMxOQDs+s+WYH8ajb15/G8iMM0L2PZPoNMvY1jb909fhlxUUhb+hksgdmjOG3mK+At3hq7+GStYz",
	"9A7TnYO7CCnqkcJUA2wsd67j91JqfNfKPGqGaZqZov1W7/buids1LRXjYeqkF4BYwetPb8RCpZV67R4R",
	"6Lug8oryVFH7Xlxi7A7jaZcpzk3S5TF+igslSDsWIAEK8IDVge3y5dlOY7r22o95N0c10LoXrr+nXZbZ",
	"Li1utgy5U4xlZgTWGgtrRwTe/+qo8tmmkLl9rPpAtlaHQdF6/clZF0aR858L8p80iT0iT4tM10N0hBLU",
	"94cocBrYbTfPk+ZdlHmxIWuI73LJVhzna6udc0xTtnFP2TZfV6rqusZrt1ldwhCYpzT1rrbrzayIdtd+",
	"hSqi6w2qxVQjC9dDs7ghy+/W1lsk97HWoeUptKpc2mN6YvPaFbGPPIkSROUT9oNAH68uFYeJq8Wkl+t4",
	"JSL2zXpUVertmb2hd3hUMEyveD5nW+h9/Ii3TakSgfSL73CKKs9gi83X3kSK8noLuSOz6ePPRUHSh96Y",
	"IHf8TK82x7WznuvP3+0+FDbKYu+MkmatOTOhMmsKM2bg4d1ODUB1U4KnPmA4kqco9gwO0UE8jgXWgwyb",
	"iUfeQzRWvtYhqDNiY8UdDyTESBp7b342teSk5YGpJErbL90lQHJjP5bG4QYk1m7UynF4fWEG28e4vZKl",
	"cyMsCBqPMQWrb+Wx7bkVVctmFBDjaMM4IC/l268TISIVNwYykcb+CsUN1Cq/DX3+3tTSaYZNWwu/dCPW",
	"nqeqhYcaK7FWgEOzVffw4/XZlXeYvOoWcYr+LO91tHWGycZjGE1GYIJ3Z15Pnb751FDBwIv4poC0V4HI",
	"jxUu3Dl/RHR6H5hXIM3knrJqnVSa2d7ltfdAw4Duixyfag9RlWAYZln6hanHs6zeeHpTDac/SHpq6r3o",
	"MQ4UH93OoGiWwjpUBkWwdNuh87NiZb4GpWU1C78NOOvPHjf/xUmijMAmaeLxny8RZW5fpPuSIebe83HP",
	"wtSeW3YE6ckf9A/BXHwF4KDcpVXn7Ivwl2AdrD04TF4HT4QmHAXMdzk8hAnDS1pIq2D9PrMk7chTmH6J",
	"JIVqkn0yEtJQOsJQg8alI8yN5zx+GOZPDq8/VOaDAor2M5n9VTWaTo6PM5bgbM2EPPmv1395PVJn1EKo",
	"uTrjzT0yLqPUVGRv3FlUSyWunFlzFEeqA8cpKTvg9W0Xaqr6+QWOHj49/G8AAAD//62Z3ongkgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	pathPrefix := path.Dir(pathToFile)

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(pathPrefix, "./common.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
