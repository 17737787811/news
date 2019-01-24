
window.onload = function (ev) {
        $(".dels").click(function () {
            if(!confirm("是否确认删除？")){
                return false
            }
        })

        $("#select").change(function () {
            $("#form").submit()
        })
    }
