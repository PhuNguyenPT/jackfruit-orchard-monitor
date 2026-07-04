// =====================================================
// MKE-S13 Variant C — Shared Overrides
// Single source of truth for everything mke_s13_case_c.scad and
// mke_s13_assembly_c.scad add on top of the base config. Both files
// `include` (not `use`) this file, so they always see identical
// values -- no manual re-declaring, no drift (same reasoning as
// mke_s13_config_b.scad).
// =====================================================
include <mke_s13_config.scad>

// =====================================================
// M3 THREAD -- SHARED CUTTING PARAMETERS (variant C only)
// =====================================================
// Variant C is the base 2-hole PCB-standoff/lid-pillar layout, extended
// so the lid's own pillar carries a real internal M3 thread section
// (phase-locked to the bottom standoff's thread) instead of a plain
// clearance bore -- see mke_s13_case_c.scad's lid() for the full
// phase-alignment derivation. These three values are used by BOTH the
// bottom standoff's thread and the lid pillar's thread, so a single M3
// screw can mesh with both without cross-threading.
m3_thread_pitch = 0.5; // ISO M3 coarse-thread pitch (mm/turn)
m3_thread_d     = 3.0 + (2 * clearance_per_side); // internal-cut diameter,
                      // print-fit compensated the same diametral way as
                      // lock_pin_d/baffle_clearance (both-sides shrinkage,
                      // not a one-sided wall gap) -- currently 3.3mm at
                      // this project's 0.2mm-nozzle setup
thread_cut_undercut = 0.01; // both thread cuts start this far below their
                      // nominal origin plane -- a CGAL boolean-safety
                      // margin (same idea as overlap_eps, just sized for
                      // a thread root rather than a flat face), NOT a
                      // print-fit number. Named here rather than left as
                      // a bare -0.01 inline because the lid pillar
                      // thread's phase-alignment math in
                      // mke_s13_case_c.scad has to reference the exact
                      // same value the bottom standoff's cut uses.
