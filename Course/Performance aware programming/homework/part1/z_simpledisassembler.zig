const std = @import("std");


const reg_table = [16]u8{
    "al", "cl", "dl", "bl", "ah", "ch", "dh", "bh",
    "ax", "cx", "dx", "bx", "sp", "bp", "si", "di"
};

pub fn main() void {
    std.debug.print("Hello, {s}\n",.{"World"});
    var arena  = std.heap.ArenaAllocator.init(std.heap.page_allocator);
    defer arena.deinit();

    const allocator = arena.allocator();

    var args = try std.process.argsWithAllocator(allocator);
    defer args.deinit();
    
    _ = args.next();

    const filename = args.next() orelse {
        std.debug.print("please provide a filename\n",.{});
        return;
    };

    const file = std.fs.cwd().openFile(filename, .{}) orelse {
        std.debug.print("failed to open file: {s}\n", .{filename});
        return;
    };
    defer file.close();

    
}