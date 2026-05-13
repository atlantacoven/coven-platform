//package space.thecoven.android
//
//import org.json.JSONObject
//
//
//data class Workload(val id: String, val data: JSONObject)
//
//
////val wk = Workload(id ="asd", data= JSONObject())
//
//data class Show(
//    val meat: List<Workload>,
//    val bread: List<Workload>,
//    val adjunct: Show? = null
//)
//
//data class ShowCompareResults(
//    val added: List<Workload>,
//    val removed: List<Workload>,
//    val changed: List<Workload>
//)
//
//fun compareWorkloadList(a: List<Workload>, b: List<Workload>): ShowCompareResults {
//    val added = mutableListOf<Workload>()
//    val removed = mutableListOf<Workload>()
//    val changed = mutableListOf<Workload>()
//
//    val aa = a.associateBy { it.id } // Map<ID: Workload>
//    val bb = b.associateBy { it.id } // Map<ID: Workload>
//
//    // for all ids in aa
//    // look in b
//    for ((id, wkld) in aa) {
//        val bwkld = bb[id]
//        if (bwkld == null) {
//            removed.add(wkld)
//        } else if (bwkld.data != wkld.data) {
//            changed.add(wkld)
//        } else {
//            // no change
//        }
//    }
//
//    for (id in (bb.keys - aa.keys)) {
//        added.add(bb[id]!!)
//    }
//    return ShowCompareResults(added,removed,changed)
//}
//
//fun compareShows(a: Show, b: Show): ShowCompareResults {
//    val added = mutableListOf<Workload>()
//    val removed = mutableListOf<Workload>()
//    val changed = mutableListOf<Workload>()
//
//
//
//}