
// 弹窗小工具
function Prompt() {

    let toast = function (c) {
        const {
            msg = "",
            icon = "success",
            position = "top-end",
        } = c;

        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            position: 'top-end',
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (d) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = d;

        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer,
        })
    }

    let error = function (c) {
        const {
            msg = "",
            title = "",
            footer = "",
        } = c;

        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer,
        })
    }

    // 实现弹出窗口选择日期，并记录表单
    async function custom(c) {
        const {
            icon="",
            msg = "",
            title = "",
            showConfirmButton=true,
            showCancelButton=true,
        } = c;

        const {value: result} = await Swal.fire({
            icon:icon,
            title: title,
            html: msg,
            backdrop:false,
            focusConfirm: false,
            showCancelButton:showCancelButton,
            confirmButtonText:'查询',
            cancelButtonText: '取消',
            showConfirmButton:showConfirmButton,
            // 增加日历选择器，在打开前会进行以下程序加载日历控件
            willOpen:() => {
                if (c.willOpen !== undefined){
                    c.willOpen();
                }
            },
            // 打开之后，通过移除disabled属性实现第一时间不打开日历，当用户点击时打开
            didOpen:() => {
                if (c.didOpen !== undefined){
                    c.willOpen();
                }
            }
        })
        if (result){
            if(result.dismiss !== Swal.DismissReason.cancel){ // 如果不是点击取消按钮关闭
                if (result.value !== ""){ // 如果输入域有输入值
                    if (c.callback !== undefined){
                        c.callback(result); // 执行回调函数,传入result
                    }
                }else{
                    c.callback(false); // 如果输入域无值,执行回调函数,传入false
                }
            }else{
                c.callback(false); // 执行回调函数,传入false
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }

}

// 房型查询按钮的函数
function checkAvailability(roomID,csrfToken) {
    document.getElementById("check-availability-button").addEventListener("click",
        function () {
            let html = `
    <form id="check-availability-form" action="" method="post" novalidate class="needs-validation">
        <div class="row">
            <div class="col">
                <div class="row" id="reservation-dates-modal">
                    <div class="col">
                        <input  required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                    </div>
                    <div class="col">
                        <input  required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                    </div>
                </div>
            </div>
        </div>
    </form>
`
            attention.custom({
                msg: html,
                title: "choose ",
                willOpen: () => {
                    const elem = document.getElementById('reservation-dates-modal');
                    const rp = new DateRangePicker(elem, {
                        format: "yyyy-mm-dd",
                        showOnFocus: false,
                        minDate: new Date(),
                    })
                },

                // didOpen:() => {
                //     console.log("didOpen called");
                //     document.getElementById('start').removeAttribute('disabled');
                //     document.getElementById('end').removeAttribute('disabled');
                // },
                callback: function (result) {
                    console.log("我被调用了",roomID);
                    let form = document.getElementById("check-availability-form");
                    let formData = new FormData(form);
                    formData.append("csrf_token", csrfToken);
                    formData.append("room_id", roomID);

                    fetch('/search-availability-json', {
                        method: "post",
                        body: formData,
                    })
                        .then(response => response.json())
                        .then(data => {
                            if (data.ok) {
                                attention.custom({
                                    icon: 'success',
                                    // msg:'<p>有空余房间</p>'
                                    // + '<button href="#!" class="swal2-confirm swal2-styled">'
                                    // + '预定</button>',
                                    msg: '<p>有空余房间</p>'
                                        + '<p><a href="/book-room?id='
                                        + data.room_id
                                        + '&s='
                                        + data.start_date
                                        + '&e='
                                        + data.end_date
                                        + '" class="btn btn-success">'
                                        + '预定</a></p>',
                                    showConfirmButton: false,
                                    showCancelButton: true,
                                    cancelButtonText: '取消',
                                })
                            } else {
                                attention.error({
                                    msg: "没有空房间"
                                })
                            }

                        })
                }
            });
        })
}